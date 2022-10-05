# Release binaries to GitHub.
release:
	@echo "==> Releasing"
	@goreleaser -p 1 --rm-dist --config .goreleaser.yaml
	@echo "==> Complete"
.PHONY: release

pre-release: 
	@echo "==> Releasing to locals"
	@goreleaser release --snapshot --rm-dist
	@echo "==> Complete"
.PHONY: pre-release