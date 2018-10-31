default: build

.PHONY: test

ginoas:
	# ▶️  Build ginoas (GIN OpenAPI Specification generator)
	GO111MODULE=on go build -o bin/ginoas cmd/ginoas/*.go

build: ginoas

test:
	# ▶️  Executing tests
	GO111MODULE=on go test ./...

test-v:
	# ▶️  Executing tests (verbose)
	GO111MODULE=on go test -v ./...

vet:
	# ▶️  Running GO Vet
	GO111MODULE=on go vet ./...

clean:
	# ▶️  Cleaning the GO environment for this project (cache and test cache only)
	GO111MODULE=on go clean -cache -testcache -i ./...
	rm -rf bin

fmt: clean
	# ▶️  Formatting source code
	GO111MODULE=on go fmt ./...

check-fmt: clean
	# ▶️  Checking source formatting
	@if [ "$$(GO111MODULE=on gofmt -d .)" != "" ]; then false; else true; fi

lint:
	# ▶️  Linting
	@if [ "$$(GO111MODULE=on golint -set_exit_status ./...)" != "" ]; then false; else true; fi

check: clean check-fmt lint vet build test

ci-trigger: clean check-fmt build vet test-v
