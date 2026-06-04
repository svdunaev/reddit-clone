.DEFAULT_GOAL: help

# mockgen is run via `go run` so the version is pinned by go.mod (no global install needed).
# Override with `make mocks MOCKGEN=mockgen` to use a binary on your PATH.
MOCKGEN   ?= go run go.uber.org/mock/mockgen
MOCKS_DIR := internal/mocks
MOCK_PKG  := mocks
# Third-party interfaces mocked in program (reflect) mode: "import/path:Interface".
# These live outside the project so they can't be auto-discovered from source.
EXTERNAL_MOCKS := k8s.io/utils/clock:Clock

.PHONY: help
help: ## shows available commands
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: cmd/server/main.go ## builds project and runs the binary file
	@go run cmd/server/main.go 

test: ## runs tests
	@echo Testing...

lint: ## runs linter
	@echo Linting...

build: cmd/server/main.go ## builds project
	@mkdir -p bin
	@go build -o bin/main cmd/server/main.go
	@./bin/main

.PHONY: generate-mocks
generate-mocks:
	@find . -name "deps.go" -type f | while read -r file; do \
		dir="$$(dirname "$$file")"; \
		mock_file="$$dir/deps_mock.go"; \
		echo "Generating mock for $$file -> $$mock_file"; \
		rm -f "$$mock_file"; \
		mockgen -source="$$file" -destination="$$mock_file" -package="$$(basename "$$dir")"; \
	done
	@echo "All mocks generated successfully."
