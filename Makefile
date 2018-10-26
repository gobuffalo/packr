TAGS ?= "sqlite"
GO_BIN ?= go

install: deps
	echo "installing packr v1"
	packr
	$(GO_BIN) install -v ./packr

tidy:
ifeq ($(GO111MODULE),on)
	$(GO_BIN) mod tidy
else
	echo skipping go mod tidy
endif

deps:
	rm -rf packrd
	rm -rf v2/packrd
	$(GO_BIN) get github.com/gobuffalo/release
	$(GO_BIN) get -tags ${TAGS} -t ./...
	$(GO_BIN) install -v ./packr
	packr clean
	make tidy

build: deps
	packr
	$(GO_BIN) build -v .
	make tidy

test:
	packr clean
	$(GO_BIN) test -tags ${TAGS} ./...
	packr clean

ci-test: deps
	$(GO_BIN) test -tags ${TAGS} -race ./...
	make tidy
	cd ./v2 && make ci-test

lint:
	gometalinter --vendor ./... --deadline=1m --skip=internal

update:
	$(GO_BIN) get -u -tags ${TAGS}
	make tidy
	packr
	make test
	make install
	make tidy
