#!/usr/bin/make -f

install: go.sum
		export GOFLAGS='-buildmode=pie'
		export CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2"
		export CGO_LDFLAGS="-Wl,-z,relro,-z,now -fstack-protector"
		go install $(BUILD_FLAGS) ./cmd/interchain-security-pd

# run all tests: unit, e2e, diff, and integration
test: 
	go test ./...

# run e2e and unit tests
test-short:
	go test ./x/... ./app/...

# run all tests with caching disabled
test-no-cache:
	go test ./... -count=1

BUILD_TARGETS := build

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/
