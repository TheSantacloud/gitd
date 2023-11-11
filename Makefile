format:
	go fmt ./...

build:
	go build -o ~/.local/bin/mgtd 

.DEFAULT_GOAL := all
all: format build 
