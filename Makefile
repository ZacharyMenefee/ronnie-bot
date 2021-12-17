.PHONY: run
run:
	ENVIRONMENT=DEV go run main.go

.PHONY: build
build:
	go build main.go
