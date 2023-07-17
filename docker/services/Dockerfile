FROM golang:1.19-alpine

RUN mkdir -p /go/src/github.com/iliyanmotovski/notification

WORKDIR /go/src/github.com/iliyanmotovski/notification

ADD services/docker-entrypoint.sh /usr/bin/docker-entrypoint
RUN chmod +x /usr/bin/docker-entrypoint

ENTRYPOINT ["docker-entrypoint"]
