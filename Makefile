# Run tests
test: gen-rest
	@echo "==> Running tests"
	@go test -cover ./... | grep -e "^[^?].*"
	@echo "==> Complete"
.PHONY: test

# Release binaries to GitHub
release:
	@echo "==> Releasing"
	@goreleaser -p 1 --rm-dist --config .goreleaser.yaml
	@echo "==> Complete"
.PHONY: release

# Pre release to debug locally
pre-release: 
	@echo "==> Releasing to locals"
	@goreleaser release --snapshot --rm-dist
	@echo "==> Complete"
.PHONY: pre-release

# Run tests with hotreload
watch-tests:
	@watch -n 2 make test
.PHONY: watch-tests

# Generate rest stub
gen-rest:
	@echo "==> Generating rest server stub"
	@oapi-codegen --config ./src/web/server/.codegen.server.yaml  ./src/web/server/.openapi.yaml
.PHONY: gen-rest

# Run deamon with hot reload
hot-deamon:
	@air -build.cmd "go build -o ./tmp/main ./src/bin/daemon/main.go"
.PHONY: hot-deamon

# Generate frogdb frontend api sdk
frontend-sdk-gen:
	@docker run --rm \
		-v ${PWD}:/local openapitools/openapi-generator-cli generate \
		-i /local/src/web/server/.openapi.yaml \
		-g typescript-axios \
		-o /local/frontend/src/apiCodegen

	# delete all files but with .ts extension
	@cd ./frontend/src/apiCodegen/; find . -not -name '*ts' -print0 | xargs -0  -I {} rm -r {}
.PHONY: frontend-sdk-gen