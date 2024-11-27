BIN           := $(PWD)/_bin
CACHE         := $(PWD)/_cache
GOPATH        := $(CACHE)/go
PATH          := $(BIN):$(PATH)
SHELL         := env PATH=$(PATH) GOPATH=$(GOPATH) /bin/sh
PROVIDER_NAME := terraform-provider-github

# Versions
go_version           := 1.23.3
golangci_version     := 1.62.2
tfplugindocs_version := 0.20.1

# Operating System and Architecture
os ?= $(shell uname|tr A-Z a-z)

ifeq ($(shell uname -m),x86_64)
  arch   ?= amd64
endif
ifeq ($(shell uname -m),arm64)
  arch   ?= arm64
endif

.PHONY: all
all: format lint install docs test

.PHONY: tools
tools: $(BIN)/go $(BIN)/golangci-lint $(GOPATH)/bin/tfplugindocs

# Setup Go
go_package_name := go$(go_version).$(os)-$(arch)
go_package_url  := https://go.dev/dl/$(go_package_name).tar.gz
go_install_path := $(BIN)/go-$(go_version)-$(os)-$(arch)

$(BIN)/go:
	@mkdir -p $(BIN)
	@mkdir -p $(GOPATH)
	@echo "Downloading Go $(go_version) to $(go_install_path)..."
	@curl --silent --show-error --fail --create-dirs --output-dir $(BIN) -O -L $(go_package_url)
	@tar -C $(BIN) -xzf $(BIN)/$(go_package_name).tar.gz && rm $(BIN)/$(go_package_name).tar.gz
	@mv $(BIN)/go $(go_install_path)
	@ln -s $(go_install_path)/bin/go $(BIN)/go

# Setup golangci
golangci_package_name := golangci-lint-$(golangci_version)-$(os)-$(arch)
golangci_package_url  := https://github.com/golangci/golangci-lint/releases/download/v$(golangci_version)/$(golangci_package_name).tar.gz
golangci_install_path := $(BIN)/$(golangci_package_name)

$(BIN)/golangci-lint:
	@mkdir -p $(BIN)
	@echo "Downloading golangci-lint $(golangci_version) to $(BIN)/golangci-lint-$(golangci_version)..."
	@curl --silent --show-error --fail --create-dirs --output-dir $(BIN) -O -L $(golangci_package_url)
	@tar -C $(BIN) -xzf $(BIN)/$(golangci_package_name).tar.gz && rm $(BIN)/$(golangci_package_name).tar.gz
	@ln -s $(golangci_install_path)/golangci-lint $(BIN)/golangci-lint

# Setup tfplugindocs
$(GOPATH)/bin/tfplugindocs: $(BIN)/go
	@go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v$(tfplugindocs_version)

.PHONY: update
update: $(BIN)/go
	@echo "Updating dependencies..."
	@go get -u
	@go mod tidy

.PHONY: build
build: update
	@echo "Building..."
	@go build ./...

.PHONY: install
install: update
	@echo "Installing provider..."
	@go install ./...

.PHONY: format
format: tools
	@echo "Formatting..."
	@go fmt ./...

.PHONY: lint
lint: tools update
	@echo "Linting..."
	@golangci-lint run ./...

.PHONY: docs
docs: tools update install
	@echo "Generating Docs..."
	@$(GOPATH)/bin/./tfplugindocs generate -rendered-provider-name "GitHub" >/dev/null

.PHONY: test
test: install
	@echo "Testing..."
	@cd internal/provider && TF_ACC=1 go test -count=1 -v

.PHONY: clean
clean:
	@echo "Removing the $(CACHE) directory..."
	@go clean -modcache
	@rm -rf $(CACHE)
	@echo "Removing the $(BIN) directory..."
	@rm -rf $(BIN)
