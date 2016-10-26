FROM alpine
MAINTAINER Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
RUN apk add -U ca-certificates && rm -rf /var/cache/apk/*
ADD website /
CMD ["/website", "/config.ini"]
EXPOSE 8080
