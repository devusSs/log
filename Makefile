checkinstall:
	@echo "=====> Checking Go installation..."
	@go version || (echo "Go is not installed. Please install go and try again." && exit 1)

deps: checkinstall
	@echo "=====> Checking dependencies"
	@mkdir -p $(shell go env GOPATH)/bin
	@if [ ! -x "$(shell go env GOPATH)/bin/golangci-lint" ]; then \
	    echo "Installing golangci-lint" && \
	    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin; \
	fi
	@if [ ! -x "$(shell go env GOPATH)/bin/gocritic" ]; then \
	    echo "Installing gocritic" && \
	    go install -v github.com/go-critic/go-critic/cmd/gocritic@latest; \
	fi
	@if [ ! -x "$(shell go env GOPATH)/bin/goimports" ]; then \
	    echo "Installing goimports" && \
	    go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@if [ ! -x "$(shell go env GOPATH)/bin/golines" ]; then \
	    echo "Installing golines" && \
	    go install github.com/segmentio/golines@latest; \
	fi

clean: deps
	@echo "=====> Cleaning up"
	@go mod tidy || exit 1
	@go vet ./... || exit 1
	@go fmt ./... || exit 1
	@golines -w ./ || exit 1
	@goimports -l ./ || exit 1
	@golangci-lint run ./... || exit 1
	@gocritic check ./... || exit 1

test: clean
	@echo "=====> Testing"
	@go clean -testcache && go test -cover -v -coverprofile=coverage.out ./... || exit 1
	@go tool cover -html=coverage.out || exit 1
	@rm coverage.out || exit 1
	@echo "=====> Testing passed, results in browser"

.PHONY: checkinstall deps clean test

.DEFAULT_GOAL := test