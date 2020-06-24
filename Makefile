Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"

run: build 
	./build/debug/adanos-alert --enable_migrate

run-only:
	./build/debug/adanos-alert --enable_migrate --listen :19998

run-agent: build-agent
	./build/debug/adanos-agent

run-proxy: build-proxy
	cat go.mod | build/debug/adanos-proxy --adanos-server http://localhost:29999 --adanos-server http://localhost:19999 --tag test --tag local --meta "abc=def" --meta "hello=world" --meta "fine" --origin cli

run-dashboard:
	cd dashboard && npm run serve

build-dashboard:
	cd dashboard && yarn build

build-agent:
	go build -race -ldflags $(LDFLAGS) -o build/debug/adanos-agent cmd/agent/main.go

build:
	go build -race -ldflags $(LDFLAGS) -o build/debug/adanos-alert cmd/server/main.go
	cp api/view/*.html build/debug/

build-proxy:
	go build -race -ldflags $(LDFLAGS) -o build/debug/adanos-proxy cmd/proxy/main.go

build-deploy-release: static-gen
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o .ansible/roles/server/files/adanos-alert-server cmd/server/main.go

deploy-server: build-deploy-release
	cd .ansible && ansible-playbook -i hosts playbook.yml --limit adanos-alert-server-prod

build-release:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/adanos-alert-darwin cmd/server/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/adanos-alert.exe cmd/server/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/release/adanos-alert-linux cmd/server/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags $(LDFLAGS) -o build/release/adanos-alert-arm cmd/server/main.go

static-gen: build-dashboard
	esc -pkg api -o api/static.go -prefix=dashboard/dist dashboard/dist

proto-build:
	protoc --go_out=plugins=grpc:. rpc/protocol/*.proto

doc-gen:
	swag init -g api/provider.go

clean:
	rm -fr build/debug/adanos-alert build/release/adanos-alert*

.PHONY: run build build-release clean build-dashboard run-dashboard static-gen doc-gen proto-build build-release-linux
