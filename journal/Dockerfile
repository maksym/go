FROM golang:1.8.0
RUN apt-get update
RUN apt-get install -y libsystemd-dev
RUN go get github.com/coreos/go-systemd/sdjournal
ADD journal.go /journal/journal.go
ADD start /journal/start
WORKDIR /journal
RUN go build
VOLUME /journal/log
VOLUME /journal/tail
CMD ["./start"]
