FROM golang:1.19-alpine

RUN mkdir -p /go/src/github.com/iliyanm/notification

WORKDIR /go/src/github.com/iliyanm/notification

ADD services/docker-entrypoint.sh /usr/bin/docker-entrypoint
RUN chmod +x /usr/bin/docker-entrypoint

ENTRYPOINT ["docker-entrypoint"]
