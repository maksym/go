FROM golang:1.8.0
RUN go get github.com/coreos/go-systemd/sdjournal
ADD journal.go /journal/journal.go
WORKDIR /journal
RUN go build
VOLUME /journal/log
VOLUME /journal/tail
ENTRYPOINT ["./journal", "/journal/log", "/journal/tail"]
