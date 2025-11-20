BIN           := $(PWD)/.local/bin
CACHE         := $(PWD)/.local/cache
GOPATH        := $(CACHE)/go
PATH          := $(BIN):$(PATH)
SHELL         := env PATH=$(PATH) GOPATH=$(GOPATH) /bin/sh
PROVIDER_NAME := terraform-provider-github

# Versions
go_version           := 1.25.1
golangci_version     := 2.4.0
tfplugindocs_version := 0.22.0
actionlint_version   := 1.7.7
shellcheck_version   := 0.11.0

# Operating System and Architecture
os ?= $(shell uname|tr A-Z a-z)

ifeq ($(shell uname -m),x86_64)
  arch   ?= amd64
endif
ifeq ($(shell uname -m),arm64)
  arch     ?= arm64
  arch_alt ?= aarch64
endif

.PHONY: all
all: format lint install docs test

.PHONY: tools
tools: $(BIN)/go $(BIN)/golangci-lint $(BIN)/tfplugindocs $(BIN)/actionlint $(BIN)/shellcheck

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
	@echo "Downloading golangci-lint $(golangci_version) to $(BIN)/golangci-lint-$(golangci_version)..." #TODO: Update this line to use golangci_install_path
	@curl --silent --show-error --fail --create-dirs --output-dir $(BIN) -O -L $(golangci_package_url)
	@tar -C $(BIN) -xzf $(BIN)/$(golangci_package_name).tar.gz && rm $(BIN)/$(golangci_package_name).tar.gz
	@ln -s $(golangci_install_path)/golangci-lint $(BIN)/golangci-lint

# Setup tfplugindocs
tfplugindocs_package_name := tfplugindocs_$(tfplugindocs_version)_$(os)_$(arch)
tfplugindocs_package_url  := https://github.com/hashicorp/terraform-plugin-docs/releases/download/v$(tfplugindocs_version)/$(tfplugindocs_package_name).zip
tfplugindocs_install_path := $(BIN)/$(tfplugindocs_package_name)

$(BIN)/tfplugindocs:
	@mkdir -p $(BIN)
	@echo "Downloading tfplugindocs $(tfplugindocs_version) to $(tfplugindocs_install_path)..."
	@mkdir -p $(tfplugindocs_install_path) # actionlint isn't packaged in a directory
	@curl --silent --show-error --fail --create-dirs --output-dir $(tfplugindocs_install_path) -O -L $(tfplugindocs_package_url)
	@cd $(tfplugindocs_install_path) && unzip $(tfplugindocs_package_name).zip && rm $(tfplugindocs_package_name).zip && cd -
	@ln -s $(tfplugindocs_install_path)/tfplugindocs $(BIN)/tfplugindocs

# Setup actionlint
actionlint_package_name := actionlint_$(actionlint_version)_$(os)_$(arch)
actionlint_package_url  := https://github.com/rhysd/actionlint/releases/download/v$(actionlint_version)/$(actionlint_package_name).tar.gz
actionlint_install_path := $(BIN)/$(actionlint_package_name)

$(BIN)/actionlint:
	@mkdir -p $(BIN)
	@echo "Downloading actionlint $(actionlint_version) to $(BIN)/actionlint-$(actionlint_version)..."
	@curl --silent --show-error --fail --create-dirs --output-dir $(BIN) -O -L $(actionlint_package_url)
	@mkdir -p $(actionlint_install_path) # actionlint isn't packaged in a directory
	@tar -C $(actionlint_install_path) -xzf $(BIN)/$(actionlint_package_name).tar.gz && rm $(BIN)/$(actionlint_package_name).tar.gz
	@ln -s $(actionlint_install_path)/actionlint $(BIN)/actionlint

# Setup shellcheck
shellcheck_package_name := shellcheck-v$(shellcheck_version).$(os).$(arch_alt)
shellcheck_package_url  := https://github.com/koalaman/shellcheck/releases/download/v$(shellcheck_version)/$(shellcheck_package_name).tar.xz
shellcheck_install_path := $(BIN)/shellcheck-v$(shellcheck_version)

$(BIN)/shellcheck:
	@mkdir -p $(BIN)
	@echo "Downloading shellcheck $(shellcheck_version) to $(BIN)/shellcheck-$(shellcheck_version)..."
	@curl --silent --show-error --fail --create-dirs --output-dir $(BIN) -O -L $(shellcheck_package_url)
	@tar -C $(BIN) -xf $(BIN)/$(shellcheck_package_name).tar.xz && rm $(BIN)/$(shellcheck_package_name).tar.xz
	@ln -s $(shellcheck_install_path)/shellcheck $(BIN)/shellcheck

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
	@actionlint

.PHONY: docs
docs: tools update install
	@echo "Generating Docs..."
	@$(BIN)/./tfplugindocs generate -rendered-provider-name "GitHub" >/dev/null

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
	@echo "Removing the $(PWD)/.local directory..."
	@if [ -d "$(PWD)/.local" ]; then rmdir "$(PWD)/.local"; fi
