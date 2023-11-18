format:
	go fmt ./...

build:
	go build -o ~/.local/bin/gitd

.DEFAULT_GOAL := all
all: format build 
