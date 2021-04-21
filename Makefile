PROJECT = Gamechart
VERSION = $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))
DIR = $(shell basename $(CURDIR))

export NOCACHE := $(if $(NOCACHE),"--no-cache")
export GO111MODULE := on

ifeq ($(shell test -f .env && echo yes), yes)
	include .env
	export $(shell sed 's/=.*//' .env)
endif

## Init
init: pre-install dep

pre-install:
	@echo "Installing git hooks..."
	@git config --global core.hooksPath ./git-hooks
	@echo "Installing golangci-lint..."
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint

## Development
run-server:
	@echo "Running server..."
	go run -ldflags "-s -w -X ${DIR}/version.VERSION=${VERSION} -X ${DIR}/version.CommitID=${COMMIT_ID}" cmd/main.go -cmd server

dep:
	@echo "Loading dependencies..."
	go mod tidy

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --timeout 3m

test-short:
