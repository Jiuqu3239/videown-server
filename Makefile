GOFMT=gofmt
GC=go build

VERSION := $(shell git describe --always --tags --long)
BUILD_NODE_PAR = -ldflags "-X main.Version=$(VERSION)" #

ARCH=$(shell uname -m)
SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)

videown-server: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o videown-server main.go



videown-server-cross: videown-server-windows videown-server-linux videown-server-darwin

videown-server-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o videown-server-windows-amd64.exe main.go

videown-server-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o videown-server-linux-amd64 main.go

videown-server-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o videown-server-darwin-amd64 main.go

tools-cross: tools-windows tools-linux tools-darwin

format:
	$(GOFMT) -w main.go

clean:
	rm -rf *.8 *.o *.out *.6 *exe
	rm -rf videown-server videown-server-*
