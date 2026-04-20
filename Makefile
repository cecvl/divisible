.PHONY: run build fmt

run:
	go run ./cmd/game

build:
	mkdir -p bin
	go build -o bin/game ./cmd/game

fmt:
	gofmt -w $$(find ./cmd ./pkg -name '*.go')