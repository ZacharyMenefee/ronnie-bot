.PHONY: run
run:
	ENVIRONMENT=DEV go run ./...

.PHONY: build
build:
	go build ./...
