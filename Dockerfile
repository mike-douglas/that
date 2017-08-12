FROM golang:1.8-alpine

WORKDIR /go/src/github.com/mike-douglas

RUN apk add --no-cache git
RUN go get -v github.com/golang/dep/cmd/dep

RUN git clone -b golang http://github.com/mike-douglas/that.git

WORKDIR /go/src/github.com/mike-douglas/that

RUN cd /go/src/github.com/mike-douglas/that && dep ensure

CMD ["go", "run", "that/that.go"]
