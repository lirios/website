FROM alpine
MAINTAINER Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
ADD website /
CMD ["/website", "/config.ini"]
EXPOSE 8080
