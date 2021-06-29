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

SOURCE=cmd/gobase/main.go
WIRE_SOURCE=internal/gobase/server
PB_SOURCE=protoc

buildPb = protoc --experimental_allow_proto3_optional -I .\
 			--go_out=paths=source_relative:. \
 			--go-grpc_out=paths=source_relative:. \
 			--grpc-gateway_out=. \
 			gobase/v1/*.proto

############################################################
# 配置信息
############################################################
BINARY_NAME=baseSrv
BINARY_UNIX=$(BINARY_NAME)

.PHONY: env clean

all: service

env:
	export GO111MODULE=on GOPROXY=https://goproxy.cn,direct

service: env
	$(GOBUILD) $(STATIC) -o $(BINARY_NAME) -v $(SOURCE)

service-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(STATIC) -o $(BINARY_NAME) -v $(SOURCE)

clean:
	rm -f $(BINARY_NAME)

pb:
	cd $(PB_SOURCE) && $(buildPb)

wire:
	cd $(WIRE_SOURCE) && wire

run: service
	./$(BINARY_NAME)

deps:
	$(GOCMD) mod tidy