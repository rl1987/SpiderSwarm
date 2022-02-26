FROM golang:1.17.7-buster AS build

WORKDIR /go/src/spiderswarm
COPY . .

# Building statically linked binary. See: https://www.arp242.net/static-go.html
RUN go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,osusergo,netgo

FROM scratch
COPY --from=build /go/src/spiderswarm/spiderswarm /bin/spiderswarm

WORKDIR /tmp
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/cacert.pem
ENTRYPOINT ["/bin/spiderswarm"]
