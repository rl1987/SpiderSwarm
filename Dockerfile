FROM golang:1.16-alpine

WORKDIR /go/src/spiderswarm
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["spiderswarm"]
