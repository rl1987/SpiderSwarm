FROM golang:1.16-buster AS build

WORKDIR /go/src/spiderswarm
COPY . .

# Building statically linked binary. See: https://www.arp242.net/static-go.html
RUN go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,osusergo,netgo

FROM debian:buster-slim
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
COPY --from=build /go/src/spiderswarm/spiderswarm /bin/spiderswarm

ENTRYPOINT ["/bin/spiderswarm"]
