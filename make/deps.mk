##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(PROJECT_DIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest
YQ ?= $(LOCALBIN)/yq

## Tool Versions
KUSTOMIZE_VERSION ?= v4.5.7
CONTROLLER_TOOLS_VERSION ?= v0.11.3
CERTMANAGER_VERSION ?= v1.11.1
YQ_VERSION ?= v4

