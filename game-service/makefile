PROJECT_NAME := game-service
PKG_LIST := $(shell go list ./... | grep -v cmd)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)
ESCAPE_LIST := $(shell go list ./... | grep -v api)
FMT_LIST := $(shell go list ./... | grep -v cmd | grep -v restapi | grep -v operations)
REPOSITORY_URI := 713614461671.dkr.ecr.eu-central-1.amazonaws.com/${PROJECT_NAME}
BUILD_VAR := git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts.SemBuild
DATE_VAR := git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts.SemDate
BUILD_DATE := $(shell date -u --iso-8601=seconds)
TARGET_VAR := git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/hashes.GameService
GIT_HASH := $(shell git rev-parse --verify HEAD)

.PHONY: all lint fmt libs dep test bench race msan coverage coverhtml bins docker clean help

all: test

swagger: ## build swagger API code
	@scripts/build-api.sh

lint: ## lint the input files
	@golint -set_exit_status ${PKG_LIST}

fmt: ${GO_FILES} ## go fmt the input files
	@go fmt ${FMT_LIST}

libs: ## update to the latest modules
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git@latest
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git@latest
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git@latest
	@go get git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git@latest
	@go mod tidy

dep: fmt ## update dependencies
	@go mod tidy

test: dep ## run unit tests (default target)
	@go test -cover ${PKG_LIST}

bench: dep ## run benchmarks
	@go test -bench . -benchmem ${PKG_LIST}

race: dep ## run race detector
	@go test -race ${PKG_LIST}

msan: dep ## run memory sanitizer
	@go test -msan ${PKG_LIST}

coverage: dep ## get unit test coverage
	@scripts/coverage.sh

coverhtml: dep ## pretty print unit test coverage to browser
	@scripts/coverage.sh html
	@xdg-open cover.html

escape: dep ## run escape analysis
	@go build -tags DEBUG -gcflags="-m" ${ESCAPE_LIST}

debug: dep ## build all binaries in DEBUG mode
	$(eval BUILD_NR=$(shell scripts/increase_build.sh))
	@go build -ldflags "-X ${TARGET_VAR}=${GIT_HASH} -X ${BUILD_VAR}=${BUILD_NR} -X ${DATE_VAR}=${BUILD_DATE}" -tags DEBUG -o bin/${PROJECT_NAME} cmd/*.go

prod: dep ## build all binaries in PROD mode
	$(eval BUILD_NR=$(shell scripts/increase_build.sh))
	@go build -ldflags "-s -w -X ${TARGET_VAR}=${GIT_HASH} -X ${BUILD_VAR}=${BUILD_NR} -X ${DATE_VAR}=${BUILD_DATE}" -o bin/${PROJECT_NAME} cmd/*.go

docker: dep ## build alpine docker image
	@docker build -t ${REPOSITORY_URI}:latest -t game-service .

clean: ## clean-up the build
	@rm -f bin/*test

help:
	@echo 'targets:'
	@echo ''
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ''
