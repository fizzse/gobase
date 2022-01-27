##############################
# makefile
##############################

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
STATIC=-ldflags '-extldflags "-static"'


## wire di 文件
WIRE_SOURCE=internal/gobase/server

## pb 配置
PB_BASE_DIR=protoc
PB_FILE_API=gobase/v1/*.proto
#PB_FILE_ERRCODE=error/v1/*.proto
PB_GO_FILES=$(shell find protoc -name *.go)

BUILD_PB_CMD_V = protoc --experimental_allow_proto3_optional -I .\
 			--go_out=paths=source_relative:. \
 			--go-grpc_out=paths=source_relative:. \
 			--grpc-gateway_out=. \
 			--validate_out=paths=source_relative,lang=go:.

BUILD_PB_CMD = protoc --experimental_allow_proto3_optional -I .\
 			--go_out=paths=source_relative:. \
 			--go-grpc_out=paths=source_relative:. \
 			--grpc-gateway_out=.

############################################################
# 配置信息
############################################################
## main函数文件
SOURCE=cmd/gobase/main.go

BINARY_NAME=srv ## you server name
BINARY_UNIX=$(BINARY_NAME)

.PHONY: env clean

all: service

env:
	export GO111MODULE=on GOPROXY=https://goproxy.cn,direct

service: env
	$(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE)

service-static: env
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(STATIC) -o $(BINARY_NAME) -v $(SOURCE)

clean:
	rm -f $(BINARY_NAME)

pb:
	cd $(PB_BASE_DIR) && $(BUILD_PB_CMD) $(PB_FILE_API)

pb-clean:
	rm $(PB_GO_FILES) -rf

wire:
	cd $(WIRE_SOURCE) && wire

run: service
	./$(BINARY_NAME)

deps:
	$(GOCMD) mod tidy

## rename
RENAME_TOOLS=tools/rename.go
MOD = github.com/fizzse/gobase
rename:
	@echo [new name]: $(MOD)
	$(GOCMD) run $(RENAME_TOOLS) $(MOD)