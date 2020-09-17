# go import path
IMPORT = $(shell head go.mod | grep module | cut -f 2 -d ' ')
# project name
PROJECT = $(shell basename $(IMPORT))
# command (executable) names
COMMANDS = $(shell cd cmd; ls)
# installation prefix
PREFIX = /usr/local

# -----------------------------------------------------------------------------

.ONESHELL:
SHELL=/bin/bash

# -----------------------------------------------------------------------------

.PHONY: all dep echo

all: test run

dep:
	go get -u golang.org/x/lint/golint

echo:
	@echo IMPORT=\"$(IMPORT)\"
	@echo PROJECT=\"$(PROJECT)\"
	@echo COMMANDS=\"$(COMMANDS)\"
	@echo PREFIX=\"$(PREFIX)\"

# -----------------------------------------------------------------------------

.PHONY: test run run-build bulid clean

test:
	golint ./...
	go vet ./...
	go test ./...

run:
	go run "cmd/$(PROJECT)/main.go"

run-build: build
	"_build/$(PROJECT)"

build:
	mkdir -p "_build"
	for command in $(COMMANDS); do
		go build -o "_build/$$command" "cmd/$$command/main.go"
	done

clean:
	rm -r "_build"

# -----------------------------------------------------------------------------

.PHONY: install uninstall

install: build
	mkdir -p "$(PREFIX)/bin"
	for command in $(COMMANDS); do
		cp "_build/$$command" "$(PREFIX)/bin"
	done

uninstall:
	for command in $(COMMANDS); do
		rm -i "$(PREFIX)/bin/$$command"
	done
