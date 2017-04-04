package main

import (
	"fmt"
	"github.com/coreos/go-systemd/sdjournal"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"
)

func main() {
	journalDir := os.Args[1]
	tailFile := os.Args[2]
	fmt.Println(journalDir)
	fmt.Println(tailFile)
	var err error
	var journal *sdjournal.Journal
	if journalDir == "" {
		journal, err = sdjournal.NewJournal()
	} else {
		fmt.Printf("using journal dir: %s", journalDir)
		journal, err = sdjournal.NewJournalFromDir(journalDir)
	}
	if err != nil {
		fmt.Printf("error opening journal: %s", err)
		os.Exit(1)
	}
	defer journal.Close()
	seeked, err := journal.Next()
	if seeked == 0 || err != nil {
		fmt.Printf("unable to seek to first item in journal")
		os.Exit(1)
	}
	journal.SeekTail()
	records := make(chan Record)
	go ReadRecords(journal, records)
	for record := range records {
		fmt.Println(record)
	}
}

func ReadRecords(journal *sdjournal.Journal, c chan<- Record) {
	record := &Record{}

	termC := MakeTerminateChannel()
	checkTerminate := func() bool {
		select {
		case <-termC:
			close(c)
			return true
		default:
			return false

		}

	}

	for {
		if checkTerminate() {
			return

		}
		err := UnmarshalRecord(journal, record)
		if err != nil {
			c <- synthRecord(
				fmt.Errorf("error unmarshalling record: %s", err),
			)
			continue

		}

		//		record.InstanceId
		c <- *record

		for {
			if checkTerminate() {
				return

			}
			seeked, err := journal.Next()
			if err != nil {
				c <- synthRecord(
					fmt.Errorf("error reading from journal: %s", err),
				)
				time.Sleep(2 * time.Second)
				continue

			}
			if seeked == 0 {
				journal.Wait(2 * time.Second)
				continue

			}
			break

		}

	}

}

func synthRecord(err error) Record {
	return Record{
		Command:  "journald-cloudwatch-logs",
		Priority: ERROR,
		Message:  err.Error(),
	}
}

type Priority int

var (
	EMERGENCY Priority = 0
	ALERT     Priority = 1
	CRITICAL  Priority = 2
	ERROR     Priority = 3
	WARNING   Priority = 4
	NOTICE    Priority = 5
	INFO      Priority = 6
	DEBUG     Priority = 7
)

var PriorityJSON = map[Priority][]byte{
	EMERGENCY: []byte("\"EMERG\""),
	ALERT:     []byte("\"ALERT\""),
	CRITICAL:  []byte("\"CRITICAL\""),
	ERROR:     []byte("\"ERROR\""),
	WARNING:   []byte("\"WARNING\""),
	NOTICE:    []byte("\"NOTICE\""),
	INFO:      []byte("\"INFO\""),
	DEBUG:     []byte("\"DEBUG\""),
}

type Record struct {
	InstanceId     string       `json:"instanceId,omitempty"`
	TimeUsec       int64        `json:"-"`
	PID            int          `json:"pid" journald:"_PID"`
	UID            int          `json:"uid" journald:"_UID"`
	GID            int          `json:"gid" journald:"_GID"`
	Command        string       `json:"cmdName,omitempty" journald:"_COMM"`
	Executable     string       `json:"exe,omitempty" journald:"_EXE"`
	CommandLine    string       `json:"cmdLine,omitempty" journald:"_CMDLINE"`
	SystemdUnit    string       `json:"systemdUnit,omitempty" journald:"_SYSTEMD_UNIT"`
	BootId         string       `json:"bootId,omitempty" journald:"_BOOT_ID"`
	MachineId      string       `json:"machineId,omitempty" journald:"_MACHINE_ID"`
	Hostname       string       `json:"hostname,omitempty" journald:"_HOSTNAME"`
	Transport      string       `json:"transport,omitempty" journald:"_TRANSPORT"`
	Priority       Priority     `json:"priority" journald:"PRIORITY"`
	Message        string       `json:"message" journald:"MESSAGE"`
	MessageId      string       `json:"messageId,omitempty" journald:"MESSAGE_ID"`
	Errno          int          `json:"machineId,omitempty" journald:"ERRNO"`
	Syslog         RecordSyslog `json:"syslog,omitempty"`
	Kernel         RecordKernel `json:"kernel,omitempty"`
	Container_Name string       `json:"containerName,omitempty" journald:"CONTAINER_NAME"`
	Container_Tag  string       `json:"containerTag,omitempty" journald:"CONTAINER_TAG"`
	Container_ID   string       `json:"containerID,omitempty" journald:"CONTAINER_ID"`
}

type RecordSyslog struct {
	Facility   int    `json:"facility,omitempty" journald:"SYSLOG_FACILITY"`
	Identifier string `json:"ident,omitempty" journald:"SYSLOG_IDENTIFIER"`
	PID        int    `json:"pid,omitempty" journald:"SYSLOG_PID"`
}

type RecordKernel struct {
	Device    string `json:"device,omitempty" journald:"_KERNEL_DEVICE"`
	Subsystem string `json:"subsystem,omitempty" journald:"_KERNEL_SUBSYSTEM"`
	SysName   string `json:"sysName,omitempty" journald:"_UDEV_SYSNAME"`
	DevNode   string `json:"devNode,omitempty" journald:"_UDEV_DEVNODE"`
}

func (p Priority) MarshalJSON() ([]byte, error) {
	return PriorityJSON[p], nil
}

func UnmarshalRecord(journal *sdjournal.Journal, to *Record) error {
	err := unmarshalRecord(journal, reflect.ValueOf(to).Elem())
	if err == nil {
		// FIXME: Should use the realtime from the log record,
		// but for some reason journal.GetRealtimeUsec always fails.
		to.TimeUsec = time.Now().Unix() * 1000

	}
	return err

}

func unmarshalRecord(journal *sdjournal.Journal, toVal reflect.Value) error {
	toType := toVal.Type()

	numField := toVal.NumField()

	// This intentionally supports only the few types we actually
	// use on the Record struct. It's not intended to be generic.

	for i := 0; i < numField; i++ {
		fieldVal := toVal.Field(i)
		fieldDef := toType.Field(i)
		fieldType := fieldDef.Type
		fieldTag := fieldDef.Tag
		fieldTypeKind := fieldType.Kind()

		if fieldTypeKind == reflect.Struct {
			// Recursively unmarshal from the same journal
			unmarshalRecord(journal, fieldVal)

		}

		jdKey := fieldTag.Get("journald")
		if jdKey == "" {
			continue

		}

		value, err := journal.GetData(jdKey)
		if err != nil || value == "" {
			fieldVal.Set(reflect.Zero(fieldType))
			continue

		}

		// The value is returned with the key and an equals sign on
		// the front, so we'll trim those off.
		value = value[len(jdKey)+1:]

		switch fieldTypeKind {
		case reflect.Int:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				// Should never happen, but not much we can do here.
				fieldVal.Set(reflect.Zero(fieldType))
				continue

			}
			fieldVal.SetInt(int64(intVal))
			break
		case reflect.String:
			fieldVal.SetString(value)
			break
		default:
			// Should never happen
			panic(fmt.Errorf("Can't unmarshal to %s", fieldType))

		}

	}

	return nil

}

func MakeTerminateChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch

}
