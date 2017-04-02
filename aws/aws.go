package main

import "os"
import "fmt"
import "github.com/hpcloud/tail"
//import "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

func main() {
  tailfile := os.Args[1]
  group := os.Args[2]
  stream := os.Args[3]
  fmt.Println(group)
  fmt.Println(stream)
  t, err := tail.TailFile(tailfile, tail.Config{Follow: true})
  if err != nil {
    fmt.Errorf("Error: %s", err)
    return
  }
  for line := range t.Lines {
    fmt.Println(line.Text)
  }
}
