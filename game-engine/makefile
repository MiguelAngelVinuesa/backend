PROJECT_NAME := game-engine
PKG := git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/$(PROJECT_NAME)
PKG_LIST := $(shell go list ./... | grep -v cmd)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

.PHONY: all lint fmt libs dep test bench race msan coverage coverhtml poker bins clean ent rngtest dieharder help

all: test

lint: ## lint the input files
	@golint -set_exit_status ${PKG_LIST}

fmt: ${GO_FILES} ## go fmt the input files
	@go fmt ${PKG_LIST}

libs: ## update to the latest modules
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git@latest
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git@latest
	@go mod tidy

dep: fmt ## update dependencies
	@go mod tidy

test: dep ## run unit tests (default target)
	@go test -trimpath -cover ${PKG_LIST}

bench: dep ## run benchmarks
	@GOMAXPROCS=1 go test -bench . --tags=debug -benchmem ${PKG_LIST}

bench-weighting: dep ## run benchmarks for Weighting
	@GOMAXPROCS=1 go test -bench Weighting --tags=debug -benchmem ${PKG_LIST}

race: dep ## run race detector
	@go test -race --tags=debug ${PKG_LIST}

msan: dep ## run memory sanitizer
	@CC=clang go test -msan --tags=debug ${PKG_LIST}

escape: dep ## run escape analysis
	@go build -gcflags="-m" ${PKG_LIST} 2>&1

coverage: dep ## get unit test coverage
	@tools/coverage.sh

coverhtml: dep ## pretty print unit test coverage to browser
	@tools/coverage.sh html
	@xdg-open cover.html

cloc: dep ## run code-line analysis
	@cloc *

poker: dep ## build poker binaries
	@go build -trimpath -o bin/poker_prng_test cmd/cards/poker/prng_test/*.go
	@go build -trimpath -o bin/poker_hand_test cmd/cards/poker/hand_test/*.go

bins: poker ## build all binaries

clean:
	@rm -f bin/*test

help:
	@echo 'targets:'
	@echo ''
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
	@echo 'requirements: apt install ent rng-tools dieharder moreutils'
	@echo ''
