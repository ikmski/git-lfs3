# meta
NAME := git-lfs3
VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

GOFILES := $(shell find . -name "*.go" -type f -not -name '*_test.go' -not -path "./vendor/*")
SOURCES := $(shell find . -name "*.go" -type f)

LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'

.PHONY: all
## all
all: build

.PHONY: setup
## setup
setup:
	go get -u github.com/Songmu/make2help/cmd/make2help

.PHONY: install-deps
## install dependencies
install-deps: setup
	go mod download

.PHONY: update-deps
## update dependencies
update-deps: setup
	go mod tidy

.PHONY: test
## run tests
test:
	go test -v

.PHONY: lint
## lint
lint:
	go vet
	for pkg in $(GOFILES); do\
		golint --set_exit_status $$pkg || exit $$?; \
	done

.PHONY: run
## run
run:
	go run $(GOFILES)

.PHONY: build
## build
build: bin/$(NAME)

bin/$(NAME): $(SOURCES)
	go build \
		-a -v \
		-tags netgo \
		-installsuffix netgo \
		-ldflags "$(LDFLAGS)" \
		-o $@

.PHONY: install
## install
install:
	go install \
		-a -v \
		-tags netgo \
		-installsuffix netgo \
		-ldflags "$(LDFLAGS)" \
		.

.PHONY: clean
## clean
clean:
	go clean -i ./...
	rm -rf bin/*

.PHONY: help
## show help
help:
	@make2help $(MAKEFILE_LIST)

