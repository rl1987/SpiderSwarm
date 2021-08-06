FROM golang:1.16-alpine AS build

WORKDIR /go/src/spiderswarm
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:10-slim
COPY --from=build /go/src/spiderswarm/spiderswarm /bin
ENTRYPOINT "/bin/spiderswarm"

