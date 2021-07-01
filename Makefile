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


##############################
# pb settings
##############################

PB_PATH =protoc
PB_FILES = gobase/v1/*.proto
BUILD_PB_CMD = protoc --experimental_allow_proto3_optional -I .\
 			--go_out=paths=source_relative:. \
 			--go-grpc_out=paths=source_relative:. \
 			--grpc-gateway_out=. 


##############################
# server setting
##############################
SOURCE=cmd/gobase/main.go
BINARY_NAME=baseSrv
BINARY_UNIX=$(BINARY_NAME)
# wire setting
WIRE_SOURCE=internal/gobase/server


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
	cd $(PB_PATH) && $(BUILD_PB_CMD) $(PB_FILES)

wire:
	cd $(WIRE_SOURCE) && wire

run: service
	./$(BINARY_NAME)

deps:
	$(GOCMD) mod tidy