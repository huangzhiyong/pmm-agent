help:                           ## Display this help message.
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
	    awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

install:                        ## Install pmm-agent binary.
	go install -v ./...

install-race:                   ## Install pmm-agent binary with race detector.
	go install -v -race ./...

test:                           ## Run tests.
	go test -v ./...

test-race:                      ## Run tests with race detector.
	go test -v -race ./...

test-cover:                     ## Run tests and collect coverage information.
	go test -v -coverprofile=cover.out -covermode=count ./...

check-license:                  ## Check that all files have the same license header.
	go run .github/check-license.go

check: install check-license    ## Run checkers and linters.
	golangci-lint run

format:                         ## Run `goimports`.
	goimports -local github.com/percona/pmm-agent -l -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")
