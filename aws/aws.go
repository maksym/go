package main

import "os"
import "fmt"
import "github.com/hpcloud/tail"
import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

func awsSequenceToken(group string, stream string) string {
params := &cloudwatchlogs.DescribeLogStreamsInput{
  LogGroupName: aws.String(group),
  LogStreamNamePrefix: aws.String(stream),}
  resp, err := svc.DescribeLogStreams(params)
  if err != nil {
    fmt.Println(err.Error())
    return ""
  }
  fmt.Println(resp)
  return ""
}

func main() {
  tailfile := os.Args[1]
  group := os.Args[2]
  stream := os.Args[3]
  fmt.Println(group)
  fmt.Println(stream)
  sess := session.Must(session.NewSession())
  svc := cloudwatchlogs.New(sess)

  awsSequenceToken()

  t, err := tail.TailFile(tailfile, tail.Config{Follow: true})
  if err != nil {
    fmt.Errorf("Error: %s", err)
    return
  }
  for line := range t.Lines {
    fmt.Println(line.Text)
  }
}

//params := &cloudwatchlogs.PutLogEventsInput{
//  LogEvents: []*cloudwatchlogs.InputLogEvent{{
//    Message: aws.String("EventMessage"),
//    Timestamp: aws.Int64(1),},},
//  LogGroupName: aws.String(group),
//  LogStreamName: aws.String(stream),
//  SequenceToken: aws.String("SequenceToken"),}
// resp, err := svc.PutLogEvents(params)
