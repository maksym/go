FROM golang:1.8.0
RUN go get github.com/hpcloud/tail
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/cloudwatchlogs
ADD aws.go /aws/aws.go
WORKDIR /aws
RUN go build
VOLUME /aws/tail
ENTRYPOINT ["./aws", "/aws/tail", "temporary", "temporary"]
