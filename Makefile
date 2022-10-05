# Release binaries to GitHub.
release:
	@echo "==> Releasing"
	@goreleaser -p 1 --rm-dist --config .goreleaser.yaml
	@echo "==> Complete"
.PHONY: release