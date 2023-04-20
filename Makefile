BUILD_NAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# zzz 工具安装 https://github.com/sohaha/zzz
# go install github.com/sohaha/zzz@latest

default: build

.PHONY: run
run:
	go build -o tmpApp
	./tmpApp

.PHONY: test
test:
	go mod tidy
	go test ./... -v

.PHONY: clean
clean:
	go clean

.PHONY: dev
dev:
	zzz w

.PHONY: build
build:test
	zzz build --os mac,win,linux -P -T

