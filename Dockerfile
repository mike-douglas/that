FROM golang:1.8-alpine

WORKDIR /go/src

RUN apk add --no-cache git
RUN go get -v github.com/golang/dep/cmd/dep

RUN git clone http://github.com/mike-douglas/that.git
WORKDIR /go/src/github.com/mike-douglas/that

RUN git checkout -b golang && dep ensure

CMD ["go", "run", "that/that.go"]
