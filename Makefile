Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

run: build 
	./build/debug/adanos-alert

run-dashboard:
	cd dashboard && npm run serve

build-dashboard:
	cd dashboard && yarn build

build:
	go build -race -ldflags $(LDFLAGS) -o build/debug/adanos-alert cmd/adanos-alert/main.go
	cp api/view/*.html build/debug/

build-release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/adanos-alert-darwin cmd/adanos-alert/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/adanos-alert.exe cmd/adanos-alert/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/adanos-alert-linux cmd/adanos-alert/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags $(LDFLAGS) -o build/release/adanos-alert-arm cmd/adanos-alert/main.go

static-gen: build-dashboard
	esc -pkg api -o api/static.go -prefix=dashboard/dist dashboard/dist

doc-gen:
	swag init -g api/provider.go

clean:
	rm -fr build/debug/adanos-alert build/release/adanos-alert*

.PHONY: run build build-release clean build-dashboard run-dashboard static-gen doc-gen
