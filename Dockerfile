FROM alpine
MAINTAINER Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
ADD website /
WORKDIR /
ENTRYPOINT /website
EXPOSE 8080
