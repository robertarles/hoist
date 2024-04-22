SHELL := /bin/bash
APP_NAME := hoist
BUILD_DIR := ./bin
SRC_DIR := ./hoist
SRC_FILES := $(wildcard $(SRC_DIR)/*.go)

## Makefile for setting up a new machine
## 	Use '##' after target to document the target, for CLI help
##   e.g. targetname: ## CLI help text for this target

help: ## Show this help.
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'	
# allow for commands that are not really building anything

.PHONY: all clean

all: darwin_amd64 darwin_arm64 linux_amd64 windows_amd64 ## make binaries for each platform

darwin_amd64: $(SRC_FILES)
	@echo "Building $@..."
	env GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)_darwin_amd64 -ldflags "-X main.Version=$(cat Version.go | grep 'const Version' | cut -d' ' -f4)" $(SRC_FILES)
darwin_arm64: $(SRC_FILES)
	@echo "Building $@..."
	env GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)_darwin_arm64 -ldflags "-X main.Version=$(cat Version.go | grep 'const Version' | cut -d' ' -f4)" $(SRC_FILES)
linux_amd64: $(SRC_FILES)
	@echo "Building $@..."
	env GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)_linux_amd64 -ldflags "-X main.Version=$(cat Version.go | grep 'const Version' | cut -d' ' -f4)" $(SRC_FILES)
windows_amd64: $(SRC_FILES)
	@echo "Building $@..."
	env GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)_windows_amd64 -ldflags "-X main.Version=$(cat Version.go | grep 'const Version' | cut -d' ' -f4)" $(SRC_FILES)

clean: ## Clean the build directory
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)/*

install: ## install the app
	@echo "Installing locally"
  go install -ldflags "-X main.Version=$(cat $(SRC_DIR)/version.go | grep 'const Version' | cut -d' ' -f4)" $(SRC_FILES)

testrun: ## run hoist against ./text/sampledirs
	cp -r ./test/sampledirs/ ./test/sampledirs.copy
	hoist ./test/sampldirs.copy
	ls -alh ./test/sampldirs.copy

testrunclean: ## clean up the sample test run
	rm -rf ./test/sampledirs.copy
