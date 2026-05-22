PROVIDER_NAME         := terraform-provider-utility-functions
BUILD_DIR             := .local/builds
PLATFORMS             := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
GOLANGCI_LINT_VERSION := v2.11.4
GOVULNCHECK_VERSION   := v1.3.0
TFPLUGINDOCS_VERSION  := v0.25.0
GO_VERSION := $(shell awk '/^go /{print $$2}' go.mod)

.PHONY: all build clean docs format install lint test testacc update

all: lint test build

# Linting

.PHONY: lint-yamllint lint-golangci-lint lint-actionlint lint-shellcheck lint-go-mod lint-govulncheck lint-tfplugindocs

lint-yamllint:
	yamllint .

lint-golangci-lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

lint-actionlint:
	actionlint

lint-shellcheck:
	find . -type f -name '*.sh' \
		-not -path './.git/*' \
		-not -path './.local/*' \
	| while IFS= read -r file; do shellcheck "$${file}"; done

lint-go-mod:
	go mod tidy
	git diff --exit-code go.mod go.sum

lint-govulncheck:
	GOTOOLCHAIN=go$(GO_VERSION) go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

lint-tfplugindocs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@$(TFPLUGINDOCS_VERSION) validate 2>/dev/null

# Go Tooling

build:
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		echo "Building $(PROVIDER_NAME)-$${os}-$${arch}"; \
		CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} \
			go build -o $(BUILD_DIR)/$(PROVIDER_NAME)-$${os}-$${arch} .; \
	done

install:
	go install ./...

format:
	go fmt ./...

lint: lint-yamllint lint-golangci-lint lint-actionlint lint-shellcheck lint-go-mod lint-govulncheck lint-tfplugindocs

docs: install
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@$(TFPLUGINDOCS_VERSION) generate \
		-rendered-provider-name "Utility Functions" >/dev/null

test:
	go test -race -count=1 ./...

testacc: install
	cd internal/provider && TF_ACC=1 go test -count=1 -v

update:
	go get -u
	go mod tidy

clean:
	rm -rf .local/
