
#build stage
FROM golang:1.17 AS builder
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /data/bin/adanos-alert-server cmd/server/main.go 

#final stage
FROM ubuntu:21.04

ENV TZ=Asia/Shanghai
RUN apt-get -y update && DEBIAN_FRONTEND="nointeractive" apt install -y tzdata ca-certificates
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /data
COPY --from=builder /data/bin/adanos-alert-server /usr/local/bin/
EXPOSE 80
EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/adanos-alert-server", "--conf", "/etc/adanos-alert-server.yaml"]
