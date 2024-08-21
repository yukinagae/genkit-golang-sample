# Default target: print this help message
.PHONY: help
.DEFAULT_GOAL := help
help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'

## genkit: Run Genkit locally
.PHONY: genkit
genkit:
	genkit start

## tidy: Tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...
