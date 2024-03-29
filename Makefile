SHELL := /bin/bash
## Makefile for setting up a new machine
## 	Use '##' after target to document the target, for CLI help
##   e.g. targetname: ## CLI help text for this target
help:               ## Show this help.
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'
# allow for commands that are not really building anything
.PHONY: buildwin buildlinux buildmacarm64 buildmacamd64 buildall
buildmacamd64: ## build a mac intel binary
	GOOS=darwin GOARCH=amd64 go build -o bin/macamd64/hoist hoist/main.go
buildmacarm64: ## build a mac arm64 binary 
	GOOS=darwin GOARCH=arm64 go build -o bin/macarm64/hoist hoist/main.go
buildlinux: ## build a linux binary
	GOOS=linux GOARCH=amd64 go build -o bin/linux/hoist hoist/main.go
buildwin: ## build a windows binary
	GOOS=windows GOARCH=amd64 go build -o bin/windows/hoist.exe hoist/main.go
buildall: buildwin buildlinux buildmacarm64 buildmacamd64 ## build all binaries
builddebug:
	go build -o bin/debug/hoist hoist/main.go

