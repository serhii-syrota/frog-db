# Release binaries to GitHub.
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

# Run tests
test: 
	@echo "==> Running tests"
	@go test -cover ./... | grep -e "^[^?].*"
	@cat test_report.txt
	@echo "==> Complete"
.PHONY: test

watch-tests:
	@watch -n 2 make test
.PHONY: watch-tests
