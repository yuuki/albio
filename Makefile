COMMIT = $$(git describe --always)
PKG = github.com/yuuki/albio
PKGS = $$(go list ./... | grep -v vendor)
CREDITS = CREDITS

all: build

.PHONY: build
build:
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" $(PKG)

.PHONY: test
test: vet
	go test -v $(PKGS)

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: patch
patch: gobump
	./script/release.sh patch

.PHONY: minor
minor: gobump
	./script/release.sh minor

.PHONY: gobump
gobump:
	go get -u github.com/motemen/gobump/cmd/gobump

.PHONY: credits
credits:
	go get github.com/go-bindata/go-bindata/...
	script/credits > $(CREDITS)
ifneq (,$(git status -s $(CREDITS)))
	go generate -x .
endif
