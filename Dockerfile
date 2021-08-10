FROM golang:1.16-buster AS build

WORKDIR /go/src/spiderswarm
COPY . .

RUN go build

FROM debian:buster-slim
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
COPY --from=build /go/src/spiderswarm/spiderswarm /bin/spiderswarm
ENTRYPOINT ["/bin/spiderswarm"]
