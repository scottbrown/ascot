.DEFAULT_GOAL: build

.PHONY: build
build:
	go build -o .build/ascot github.com/scottbrown/ascot/cmd/ascot

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: security-check
security-check: supply-chain sast

.PHONY: supply-chain
supply-chain:
	govulncheck ./...

.PHONY: sast
sast:
	gosec ./...

.PHONY: sbom
sbom:
	spdx-sbom-generator
	cyclonedx-gomod app --json=true --licenses=true > sbom.cyclonedx.json

.PHONY: test
test:
	go test ./...
