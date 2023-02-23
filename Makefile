GIT_REV ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags --abbrev=0)
LDFLAGS := "-s -w"
GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)
file    = gout-black
package := $(file)_$(GOOS)_$(GOARCH)
ext     =

ifeq ("$(GOOS)", "windows")
	ext = .exe
endif

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: default
default: help

## Build:
prepare: ## Download depencies and prepare dev env
	@pre-commit install
	@go mod download
	@go mod tidy
	@go mod vendor

build:  ## Builds the cli binary
	go build -trimpath -ldflags=$(LDFLAGS) -o ./bin/$(file) main.go

build-ci: ## Optimized build for CI
	@echo $(GOOS)/$(GOARCH)
	go build -trimpath -ldflags=$(LDFLAGS) -o ./$(file)$(ext) main.go
	tar -czf $(package).tar.gz ./$(file)$(ext) ./LICENSE

release: ## Release with a new tag. Use like this: 'VERSION=v0.0.1 make release'
	git-chglog --next-tag $(VERSION) -o CHANGELOG.md
	git add CHANGELOG.md
	git commit -m "chore: update changelog for $(VERSION)"
	git tag $(VERSION)
	git push origin main $(VERSION)

chglog: ## Generate CHANGELOG.md
	@git-chglog -o CHANGELOG.md

## Test:
coverage:  ## Run test coverage suite
	@go test ./... -coverprofile=cov.out
	@go tool cover -html=cov.out
	@rm cov.out

## Help:
.PHONE: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
