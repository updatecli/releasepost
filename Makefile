GOVERSION=$(shell go version)
export GOVERSION

DOCKER_BUILDKIT=1
export DOCKER_BUILDKIT

local_bin=./dist/releasepost_$(shell go env GOHOSTOS)_$(shell go env GOHOSTARCH)/releasepost

.PHONY: build
build:
	go build -o bin/releasepost .

.PHONY: test
test:
	./bin/releasepost --config config.example.yaml

clean: ## Clean go test cache
	go clean -testcache
	rm -Rf dist

.PHONY: build
build: ## Build updatecli as a "dirty snapshot" (no tag, no release, but all OS/arch combinations)
	goreleaser build --snapshot --clean

.PHONY: build.all
build.all: ## Build updatecli for "release" (tag or release and all OS/arch combinations)
	goreleaser --clean --skip=publish,sign

.PHONY: release ## Create a new updatecli release including packages
release: ## release generate a release
	goreleaser release --clean --timeout=2h

.PHONY: release.snapshot ## Create a new snapshot release without publishing assets
release.snapshot: ## release.snapshot generate a snapshot release but do not published it (no tag, but all OS/arch combinations)
	goreleaser release --snapshot --clean --skip=publish,sign

.PHONY: version
version: ## Run the "version" updatecli's subcommand for smoke test
	"$(local_bin)" version

.PHONY: lint
lint: ## Execute the Golang's linters on updatecli's source code
	golangci-lint run

