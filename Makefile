.DEFAULT_GOAL: help

.PHONY: help
help: ## shows available commands
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: cmd/server/main.go ## builds project and runs the binary file
	@go build cmd/server/main.go 
	@./main 

test: ## runs tests
	@echo Testing...

lint: ## runs linter
	@echo Linting...