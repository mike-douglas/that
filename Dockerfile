FROM golang:1.8-alpine

WORKDIR /go/src/that
ADD . .

RUN apk add --no-cache git
RUN go get -v github.com/codegangsta/gin \
    && go get -v github.com/golang/dep/cmd/dep

RUN dep ensure

CMD ["go", "run", "that.go"]
