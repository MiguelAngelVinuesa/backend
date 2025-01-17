PROJECT_NAME := prng
PKG := git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/$(PROJECT_NAME)
PKG_LIST := $(shell go list ./... | grep -v spew | grep -v test)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)
DATE := $(shell date -u +'%Y%m%d-%H%M%S-%Z')
ENT_FILE := output/ent/$(DATE).log
RNG_TEST_FILE := output/rngtest/$(DATE).log
DIEHARDER_FILE := output/dieharder/$(DATE).log
GIT_HASH := $(shell git rev-parse HEAD)
PLUGIN_FILE := libprng.so
PLUGIN_PATH := bin/${PLUGIN_FILE}

.PHONY: all lint fmt libs dep test bench race msan coverage coverhtml plugin bin clean ent rngtest dieharder help

all: test

lint: ## lint the input files
	@golint -set_exit_status ${PKG_LIST}

fmt: ${GO_FILES} ## go fmt the input files
	@go fmt ${PKG_LIST}

dep: fmt ## update dependencies
	@go get -v -d ./...
	@go mod tidy

test: dep ## run unit tests (default target)
	@go test -cover  ${PKG_LIST}

bench: dep ## run benchmarks
	@go test -bench . --tags=debug -benchmem ${PKG_LIST}

race: dep ## run race detector
	@go test -race --tags=debug ${PKG_LIST}

msan: dep ## run memory sanitizer
	@go test -msan --tags=debug ${PKG_LIST}

coverage: dep ## get unit test coverage
	@tools/coverage.sh

coverhtml: dep ## pretty print unit test coverage to browser
	@tools/coverage.sh html
	@xdg-open cover.html

so: dep ## build prng shared library
	@go build -buildmode=c-shared -ldflags="-s -w -X main.hash=${GIT_HASH}" -o ${PLUGIN_PATH} cmd/sharedlib/*.go
	@sha1sum ${PLUGIN_PATH}
	@sudo cp bin/libprng.so /usr/local/lib
	@sudo cp bin/libprng.h /usr/local/include

bin: so ## build binaries
	@go build -ldflags="-s -w" -o bin/prng_spew cmd/spew/*.go
	@go build -ldflags="-s -w" -o bin/test_sharedlib cmd/test/sharedlib/*.go

clean:
	@rm -Rf bin
	@rm -Rf output

ent: rand ## run "ent" test on PRNG (100m uint32)
	@mkdir -p output/ent
	@bin/prng_spew | head -c 100000000 | ent | ts '[%Y-%m-%d %H:%M:%S]' | tee $(ENT_FILE)

rngtest: rand ## run "rngtest" on PRNG (25m uint32)
	@mkdir -p output/rngtest
	@bin/prng_spew | head -c 25000000 | rngtest 2>&1 | ts '[%Y-%m-%d %H:%M:%S]' | tee $(RNG_TEST_FILE)

dieharder: rand ## run all valid "dieharder" tests on PRNG (this may take a long time!)
	@mkdir -p output/dieharder
	@bin/prng_spew | dieharder -g 200 -a -Y 1 -k 2 | ts '[%Y-%m-%d %H:%M:%S]' | tee $(DIEHARDER_FILE)

help:
	@echo 'targets:'
	@echo ''
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
	@echo 'requirements: apt install ent rng-tools dieharder moreutils'
	@echo ''
