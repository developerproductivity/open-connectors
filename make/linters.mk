##@ Linters

GOLANGCI_LINT=$(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.54.1

YAMLLINT_VERSION ?= 1.28.0

SHELLCHECK=$(LOCALBIN)/shellcheck
SHELLCHECK_VERSION ?= v0.9.0

.PHONY: lint
lint: lint-go

GO_LINT_CONCURRENCY ?= 4
GO_LINT_OUTPUT ?= colored-line-number
GO_LINT_CMD = GOFLAGS="$(GOFLAGS)" GOGC=30 GOCACHE=$(GOCACHE) $(GOLANGCI_LINT) run --concurrency=$(GO_LINT_CONCURRENCY) --out-format=$(GO_LINT_OUTPUT)

.PHONY: lint-go
lint-go: $(GOLANGCI_LINT) fmt vet ## Checks Go code
	$(GO_LINT_CMD)

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(GOLANGCI_LINT_VERSION)
