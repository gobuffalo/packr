TAGS ?= "sqlite"
GO_BIN ?= go

install: deps
	$(GO_BIN) install -v .

deps:
	$(GO_BIN) get -tags ${TAGS} -t ./...

build: deps
	$(GO_BIN) build -v .

test:
	packr clean
	$(GO_BIN) test -tags ${TAGS} ./...
	packr clean

ci-test: deps
	packr clean
	$(GO_BIN) test -tags ${TAGS} -race ./...
	packr clean

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u
	$(GO_BIN) mod tidy
	make test
	make install

release-test:
	$(GO_BIN) test -tags ${TAGS} -race ./...

release:
	$(GO_BIN) get github.com/gobuffalo/release
	release -y -f version.go
