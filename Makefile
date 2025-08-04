.PHONY: all

# requires 'protoc' compiler:
# apt-get install protobuf-compiler
#
# and protoc-gen-go:
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

all:
	protoc --go_out=. ./record/.proto
	go build -o search ./cmd/search
	go build -o mkscolldb ./cmd/mkscolldb

