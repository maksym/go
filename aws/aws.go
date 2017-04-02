package main

import "fmt"

//func main() {
//    fmt.Printf("hello, world\n")
//}


// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatchlogs/
// import "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

// https://github.com/saymedia/journald-cloudwatch-logs

// https://github.com/hpcloud/tail

import "github.com/hpcloud/tail"

func main() {
  t, err := tail.TailFile("/home/max/aaa.bbb", tail.Config{Follow: true})
  for line := range t.Lines {
    fmt.Println("line: ")
    fmt.Println(line.Text)
  }
  fmt.Println("err: ")
  fmt.Println(err)
}
