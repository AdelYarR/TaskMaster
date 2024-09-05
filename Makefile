.PHONY: build

build:
		go build -o taskmaster ./cmd/main

.DEFAULT_GOAL := build