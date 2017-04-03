package main

import "os"
import "fmt"
import "github.com/coreos/go-systemd/sdjournal"

func main() {
  logdirectory := os.Args[1]
  tailfile := os.Args[2]
  fmt.Println(logdirectory)
  fmt.Println(tailfile)
}

// https://github.com/mheese/journalbeat/blob/master/journal/follow.golang
// https://godoc.org/github.com/coreos/go-systemd/sdjournal
// https://github.com/saymedia/journald-cloudwatch-logs/blob/master/main.go
