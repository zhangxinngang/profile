# -*- coding:utf-8 -*-
.PHONY: build-proto,build
build-pro:
	go install
build:build-proto
	go install
build-proto:
	protoc --proto_path=$$GOPATH/src/github.com/gogo/protobuf/protobuf:../../:. --gogo_out=.  serial.proto

