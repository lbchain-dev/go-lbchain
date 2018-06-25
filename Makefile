# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: glbchain-dev android ios glbchain-dev-cross swarm evm all test clean
.PHONY: glbchain-dev-linux glbchain-dev-linux-386 glbchain-dev-linux-amd64 glbchain-dev-linux-mips64 glbchain-dev-linux-mips64le
.PHONY: glbchain-dev-linux-arm glbchain-dev-linux-arm-5 glbchain-dev-linux-arm-6 glbchain-dev-linux-arm-7 glbchain-dev-linux-arm64
.PHONY: glbchain-dev-darwin glbchain-dev-darwin-386 glbchain-dev-darwin-amd64
.PHONY: glbchain-dev-windows glbchain-dev-windows-386 glbchain-dev-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

glbchain-dev:
	build/env.sh go run build/ci.go install ./cmd/glbchain-dev
	@echo "Done building."
	@echo "Run \"$(GOBIN)/glbchain-dev\" to launch glbchain-dev."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/glbchain-dev.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Glbchain-dev.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

glbchain-dev-cross: glbchain-dev-linux glbchain-dev-darwin glbchain-dev-windows glbchain-dev-android glbchain-dev-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-*

glbchain-dev-linux: glbchain-dev-linux-386 glbchain-dev-linux-amd64 glbchain-dev-linux-arm glbchain-dev-linux-mips64 glbchain-dev-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-*

glbchain-dev-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/glbchain-dev
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep 386

glbchain-dev-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/glbchain-dev
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep amd64

glbchain-dev-linux-arm: glbchain-dev-linux-arm-5 glbchain-dev-linux-arm-6 glbchain-dev-linux-arm-7 glbchain-dev-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep arm

glbchain-dev-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/glbchain-dev
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep arm-5

glbchain-dev-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/glbchain-dev
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep arm-6

glbchain-dev-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/glbchain-dev
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep arm-7

glbchain-dev-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/glbchain-dev
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep arm64

glbchain-dev-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/glbchain-dev
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep mips

glbchain-dev-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/glbchain-dev
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep mipsle

glbchain-dev-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/glbchain-dev
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep mips64

glbchain-dev-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/glbchain-dev
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-linux-* | grep mips64le

glbchain-dev-darwin: glbchain-dev-darwin-386 glbchain-dev-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-darwin-*

glbchain-dev-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/glbchain-dev
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-darwin-* | grep 386

glbchain-dev-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/glbchain-dev
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-darwin-* | grep amd64

glbchain-dev-windows: glbchain-dev-windows-386 glbchain-dev-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-windows-*

glbchain-dev-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/glbchain-dev
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-windows-* | grep 386

glbchain-dev-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/glbchain-dev
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/glbchain-dev-windows-* | grep amd64
