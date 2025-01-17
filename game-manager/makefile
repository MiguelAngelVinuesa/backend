PROJECT_NAME := game-manager
PKG := git-codecommit.eu-central-1.amazonaws.com/v1/repos/${PROJECT_NAME}.git/$(PROJECT_NAME)
PKG_LIST := $(shell go list ./... | grep -v cmd | grep -v dstore)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

.PHONY: all lint fmt libs dep test bench race msan coverage coverhtml bins clean help

all: test

dswagger: ## build d-store API client code
	@scripts/build-d-store-api.sh

lint: ## lint the input files
	@golint -set_exit_status ${PKG_LIST}

fmt: ${GO_FILES} ## go fmt the input files
	@go fmt ${PKG_LIST}

libs: ## update to the latest modules
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git@latest
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git@latest
	@go mod tidy

dep: fmt ## update dependencies
	@go mod tidy

test: dep ## run unit tests (default target)
	@go test -trimpath -cover ${PKG_LIST}

bench: dep ## run benchmarks
	@go test -bench . --tags=debug -benchmem ${PKG_LIST}

race: dep ## run race detector
	@go test -race --tags=debug ${PKG_LIST}

msan: dep ## run memory sanitizer
	@go test -msan --tags=debug ${PKG_LIST}

escape: dep ## run escape analysis
	@go build -gcflags="-m" ${PKG_LIST}

coverage: dep ## get unit test coverage
	@scripts/coverage.sh

coverhtml: dep ## pretty print unit test coverage to browser
	@scripts/coverage.sh html
	@xdg-open cover.html

bins: ## build all binaries

clean:
	@rm -f bin/*test

help:
	@echo 'targets:'
	@echo ''
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
