
#build stage
FROM golang:1.17 AS builder
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /data/bin/adanos-alert-agent cmd/agent/main.go 

#final stage
FROM ubuntu:21.04
WORKDIR /data
COPY --from=builder /data/bin/adanos-alert-agent /usr/local/bin/

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/adanos-alert-agent", "--conf", "/etc/adanos-alert-agent.yaml"]
