NAME=docker-hosts

all: build

.godeps:
	gvp init
	gvp in gpm install

init: .godeps
	mkdir -p stage

build: stage/$(NAME)

stage/$(NAME): init *.go
	gvp in go build -v -o stage/$(NAME) ./...

clean:
	rm -rf stage release .godeps
