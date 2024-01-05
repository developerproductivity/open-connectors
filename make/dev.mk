##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: test
test:  fmt vet  ## Run tests.
	$(GO) test ./... -coverprofile cover.out


