PKG := "github.com/flowshot-io/navigator-cli"
PKG_LIST := $(shell go list ${PKG}/...)

.PHONY: check
.DEFAULT_GOAL := check
check: dep fmt vet test ## Check project

.PHONY: vet
vet: ## Vet the files
	@go vet ${PKG_LIST}

.PHONY: fmt
fmt: ## Style check the files
	@gofmt -s -w .

.PHONY: test
test: ## Run tests
	@go test -short ${PKG_LIST}

.PHONY: race
race: ## Run tests with data race detector
	@go test -race ${PKG_LIST}

.PHONY: benchmark
benchmark: ## Run benchmarks
	@go test -run="-" -bench=".*" ${PKG_LIST}

.PHONY: dep
dep:
	@go mod download
	@go mod tidy

.PHONY: run
run:
	@go run ./cmd/cli

.PHONY: build
build:
	@go build ./cmd/cli

.PHONY: gen-key
gen-key:
	@cosign generate-key-pair

.PHONY: gen-proto
gen-proto:
	@./scripts/gen-proto.sh navigatorservice

.PHONY: release-dry
release-dry:
	@goreleaser release --snapshot --rm-dist

.PHONY: release
release:
	@./scripts/release.sh $(version)