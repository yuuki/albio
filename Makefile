COMMIT = $$(git describe --always)
PKG = github.com/yuuki/albio
PKGS = $$(go list ./... | grep -v vendor)

all: build

.PHONY: build
build:
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" $(PKG)

.PHONY: test
test:
	go test -v $(PKGS)
