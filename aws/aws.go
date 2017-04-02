package main

import "os"
import "fmt"
import "github.com/hpcloud/tail"

func main() {
  tailfile := os.Args[1]
  t, err := tail.TailFile(tailfile, tail.Config{Follow: true})
  for line := range t.Lines {
    fmt.Println("line: ")
    fmt.Println(line.Text)
  }
  fmt.Println("err: ")
  fmt.Println(err)
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatchlogs/
// import "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

// https://github.com/saymedia/journald-cloudwatch-logs

// https://github.com/hpcloud/tail
