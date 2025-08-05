.PHONY: all

all:
	go build -o scollsrch ./cmd/search
	go build -o mkscolldb ./cmd/mkscolldb

