##@ Docs

MDBOOK_VERSION ?= v0.4.30

.PHONY: book
book:
	ID=$(shell id -u):$(shell id -g) MDBOOK_VERSION=$(MDBOOK_VERSION) \
		docker compose -f deploy/book/docker-compose.yaml run --rm --build mdbook

