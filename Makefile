.DEFAULT_TARGET: help

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run HTTP server locally on port 8080
	@go run cmd/server/main.go

.PHONY: test
test: ## Execute the tests in the development environment
	@go test ./... -count=1 -timeout 2m

.PHONY: coverage
coverage: ## Generate test coverage in the development environment
	go test ./... -coverprofile=/tmp/coverage.out -coverpkg=./...
	go tool cover -html=/tmp/coverage.out

.PHONY: lint
lint: ## Execute syntatic analysis in the code and autofix minor problems
	@golangci-lint -c .code_quality/.golangci.yml run --fix

.PHONY: ci
ci: lint test ## Execute the tests and lint commands

.PHONY: check
check: ## Check for vulnerabilities in the dependencies
	@govulncheck -version || go golang.org/x/vuln/cmd/govulncheck@latest
	@go mod tidy
	@govulncheck ./...
