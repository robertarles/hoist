SHELL := /bin/bash
## Makefile for setting up a new machine
## 	Use '##' after target to document the target, for CLI help
##   e.g. targetname: ## CLI help text for this target

help: ## Show this help.
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'	
# allow for commands that are not really building anything

buildlocal: ## build the app
	go build -o bin/hoist -ldflags "-X main.Version=$(cat Version.go | grep 'const Version' | cut -d' ' -f4)" ./hoist/main.go; 

installlocal: ## install the app
	go install -ldflags "-X main.Version=$(cat Version.go | grep 'const Version' | cut -d' ' -f4)" ./hoist/ 

testrun: ## run hoist against ./text/sampledirs
	cp -r ./test/sampledirs/ ./test/sampledirs.copy
	hoist ./test/sampldirs.copy
	ls -alh ./test/sampldirs.copy

testrunclean: ## clean up the sample test run
	rm -rf ./test/sampledirs.copy
