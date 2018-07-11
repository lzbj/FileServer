PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
BINARY := fileserver
TEST_FLAGS=-test.short -v

REV=$(shell git log --max-count=1 --pretty="format:%h" .)

GO_VER=$(shell go version|grep "go version"|cut -d' ' -f3|sed "s/[\s\t]*//")

quick: dep
	go build -ldflags "-X FileServer.REVISION=$(REV) -X FileServer.GO_VERSION=$(GO_VER)" -o $(BINARY)

PKGS=$(shell go list ./...)

fmt:
	$(foreach i,$(PKGS),go fmt $(i);)

unit:
	$(foreach i,$(PKGS),go test $(i) $(TEST_FLAGS) || exit ;)

dep: # install dep
	go get github.com/golang/dep
	dep ensure
