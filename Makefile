GO ?= go

.PHONY: all gen build test vet fmt example update-data

all: gen build test

## gen: regenerate embedded assets + routing table from the submodule
gen:
	$(GO) run ./cmd/gmd-gen

build:
	$(GO) build ./...

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

example:
	$(GO) run ./example

## update-data: pull the latest minecraft-data and regenerate
update-data:
	git submodule update --remote minecraft-data
	$(MAKE) gen
	$(MAKE) test
