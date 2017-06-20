FROM golang:1.7.3

MAINTAINER Zonesan <chaizs@asiainfo.com>

ENV TIME_ZONE=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone

COPY . /go/src/github.com/ocmanager/svc-amount

WORKDIR /go/src/github.com/ocmanager/svc-amount

EXPOSE 8080

RUN go build

ENTRYPOINT ["./svc-amount"]
