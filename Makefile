.DEFAULT_GOAL: build

.PHONY: build
build:
	go build -o .build/ascot github.com/scottbrown/ascot

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
