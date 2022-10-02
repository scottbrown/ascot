.DEFAULT_GOAL: build

.PHONY: build
build:
	go build -o .build/ascot github.com/scottbrown/ascot

.PHONY: fmt
fmt:
	go fmt ./...
