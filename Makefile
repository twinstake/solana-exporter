# golangci-lint (and its bundled gofumpt formatter) is pinned in the isolated
# tools/ module. Keep the version in sync with .github/workflows/golangci-lint.yml.
# Update with:
#   cd tools && go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@vX.Y.Z
GOLANGCI_LINT := go tool -modfile=tools/go.mod golangci-lint

.PHONY: build test fmt lint lint-fix

build:
	CGO_ENABLED=0 go build ./cmd/solana-exporter

test:
	go test ./...

# Format the codebase (gofumpt, via golangci-lint's formatters).
fmt:
	$(GOLANGCI_LINT) fmt

# Run the linter.
lint:
	$(GOLANGCI_LINT) run

# Run the linter and auto-fix what can be fixed.
lint-fix:
	$(GOLANGCI_LINT) run --fix
