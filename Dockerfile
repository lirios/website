FROM golang:1.6-alpine
MAINTAINER Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
ADD . /go/src/github.com/lirios/website
RUN rm -f website
RUN go install github.com/lirios/website
RUN rm -rf /go/src
WORKDIR /go/bin
ENTRYPOINT /go/bin/website
EXPOSE 8080
