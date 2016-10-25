FROM alpine
MAINTAINER Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
ADD website /
CMD ["/website"]
EXPOSE 8080
