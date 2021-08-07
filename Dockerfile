FROM golang:1.16-alpine AS build

WORKDIR /go/src/spiderswarm
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:10-slim
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
COPY --from=build /go/src/spiderswarm/spiderswarm /bin
ENTRYPOINT "/bin/spiderswarm"

