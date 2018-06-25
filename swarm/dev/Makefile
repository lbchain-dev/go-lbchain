.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/lbchain-devchain/go-lbchain-dev/cmd/swarm
	go build -o bin/glbchain-dev     github.com/lbchain-devchain/go-lbchain-dev/cmd/glbchain-dev
	go build -o bin/bootnode github.com/lbchain-devchain/go-lbchain-dev/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/lbchain-devchain/go-lbchain-dev/swarm/...
