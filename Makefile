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

############################################################
# 配置信息
############################################################
BINARY_NAME=mqttRoute
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

run: service
	./$(BINARY_NAME)

deps:
	$(GOCMD) mod tidy