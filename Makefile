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
	cd ./v2 && make test
	packr clean

ci-test: deps
	$(GO_BIN) test -tags ${TAGS} -race ./...

update:
	$(GO_BIN) get -u
	$(GO_BIN) mod tidy
	make test
	make install

release:
	cd ./v2 && make release
