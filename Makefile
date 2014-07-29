NAME=docker-hosts
BIN=.godeps/bin

SOURCES=main.go hosts.go

.PHONY: all init build tools release clean

all: build

.godeps:
	gvp init
	gvp in gpm install

init: .godeps
	mkdir -p stage

build: stage/$(NAME)

stage/$(NAME): init $(SOURCES)
	gvp in go build -v -o $@ ./...

$(BIN)/gpm: init
	curl -s -L -o $@ https://github.com/pote/gpm/raw/v1.2.3/bin/gpm
	chmod +x $@

$(BIN)/gvp: init
	curl -s -L -o $@ https://github.com/pote/gvp/raw/v0.1.0/bin/gvp
	chmod +x $@

tools: $(BIN)/gpm $(BIN)/gvp

release/$(NAME): tools $(SOURCES)
	docker run \
		-i -t \
		-v $(PWD):/gopath/src/app \
		-w /gopath/src/app \
		google/golang:1.3 \
		$(BIN)/gvp in go build -v -o $@ ./...

release: release/$(NAME)

docker: release
	docker build --tag=blalor/$(NAME) .
	docker push blalor/$(NAME)

clean:
	rm -rf stage release .godeps
