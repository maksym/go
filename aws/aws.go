package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/hpcloud/tail"
	"os"
	"strconv"
)

func main() {
	tailfile := os.Args[1]
	group := os.Args[2]
	stream := os.Args[3]
	fmt.Println(group)
	fmt.Println(stream)
	sess := session.Must(session.NewSession())
	svc := cloudwatchlogs.New(sess)
	describeParams := &cloudwatchlogs.DescribeLogStreamsInput{LogGroupName: aws.String(group), LogStreamNamePrefix: aws.String(stream)}
	describeResp, describeErr := svc.DescribeLogStreams(describeParams)
	if describeErr != nil {
		fmt.Println(describeErr.Error())
	}
	fmt.Println("DescribeLogStreams: ", describeResp)
	sequenceToken := describeResp.LogStreams[0].UploadSequenceToken
	fmt.Println("sequenceToken: ", sequenceToken)

	t, err := tail.TailFile(tailfile, tail.Config{Follow: true})
	if err != nil {
		fmt.Errorf("Error: %s", err)
		return
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
		var dat map[string]interface{}
		bytes := []byte(line.Text)
		if err := json.Unmarshal(bytes, &dat); err != nil {
			fmt.Errorf("Error: %s", err)
		}
		timestamp, err := strconv.ParseInt(fmt.Sprintf("%.0f", dat["timestamp"].(float64)), 10, 64)
		fmt.Println(dat)
		putParams := &cloudwatchlogs.PutLogEventsInput{
			LogEvents:     []*cloudwatchlogs.InputLogEvent{{Message: aws.String(line.Text), Timestamp: aws.Int64(timestamp)}},
			LogGroupName:  aws.String(group),
			LogStreamName: aws.String(stream),
			SequenceToken: sequenceToken}
		putResp, putErr := svc.PutLogEvents(putParams)
		if err != nil {
			fmt.Errorf("Error: %s", putErr)
			return
		}
		fmt.Println("PutLogEvents: ", putResp)
		sequenceToken = putResp.NextSequenceToken
	}
}
