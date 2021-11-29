
#build stage
FROM golang:1.17 AS builder
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /data/bin/adanos-alert-server cmd/server/main.go 

#final stage
FROM alpine:latest
WORKDIR /data
COPY --from=builder /data/bin/adanos-alert-server /usr/local/bin/adanos-alert-server
EXPOSE 80
EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/adanos-alert-server", "--conf", "/etc/adanos-alert-server.yaml"]
