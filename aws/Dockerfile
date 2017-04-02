FROM golang:1.8.0
ADD aws.go /aws/aws.go
WORKDIR /aws
RUN go get github.com/hpcloud/tail
RUN go build
VOLUME /aws/tail
ENTRYPOINT ["./aws", "/aws/tail"]
