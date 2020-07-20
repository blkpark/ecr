MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT := $(dir $(MAKEFILE))
BIN := $(PROJECT)bin
OUT := $(BIN)/ecr

CONTAINER := golang:1.14
WORKDIR := /opt/ecr

BASE := blkpark/scratch:latest
TAG := blkpark/ecr:latest

export GOFLAGS=-mod=vendor
export CGO_ENABLED=0

default: build

deps:
	go mod vendor

build:
	@go build -a -installsuffix cgo -ldflags '-w' -o ${OUT}

run:
	go run main.go

clean:
	rm -f ${OUT}

test:
	go test -v ./...

release: clean
	docker pull ${CONTAINER}
	docker run -i -v ${PROJECT}:${WORKDIR} --rm ${CONTAINER} make -C ${WORKDIR}
