SHELL := /bin/bash

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

generate: oapi-codegen
	$(OAPICODEGEN) -generate types -o pkg/api/types.gen.go -package api openapi.yml
	$(OAPICODEGEN) -generate gin -o pkg/api/server.gen.go -package api openapi.yml

# Run go fmt against code
fmt: gci
	go mod tidy
	go fmt ./...
	find . -type f -name '*.go' -exec $(GCI) write -s standard -s default -s "prefix(github.com/cheina97/timeserver)" {} \;

# Run markdownlint against markdown files
md-lint: markdownlint
	markdownlint '**/*.md' -c .markdownlint.yml

# Run golangci-lint against code
go-lint: golangci-lint
	$(GOLANGCILINT) run --new

# Install oapi-codegen if not available
oapi-codegen:
ifeq (, $(shell which oapi-codegen))
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
OAPICODEGEN=$(GOBIN)/oapi-codegen
else
OAPICODEGEN=$(shell which oapi-codegen)
endif

# Install gci if not available
gci:
ifeq (, $(shell which gci))
	@go install github.com/daixiang0/gci@v0.13.5
GCI=$(GOBIN)/gci
else
GCI=$(shell which gci)
endif

# Install golangci-lint if not available
golangci-lint:
ifeq (, $(shell which golangci-lint))
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
GOLANGCILINT=$(GOBIN)/golangci-lint
else
GOLANGCILINT=$(shell which golangci-lint)
endif

markdownlint:
ifeq (, $(shell which markdownlint))
	@echo "markdownlint is not installed. Please install it: https://github.com/igorshubovych/markdownlint-cli#installation"
	@exit 1
else
MARKDOWNLINT=$(shell which markdownlint)
endif
