FROM golang:1.5

ADD . /go/src/github.com/joshansen/WineDatabase
WORKDIR /go/src/github.com/joshansen/WineDatabase
RUN go get