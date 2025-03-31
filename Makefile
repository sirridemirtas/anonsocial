.PHONY: dev build clean run

deps:
	go mod download

dev:
	GIN_MODE=debug go run github.com/air-verse/air@latest

build: deps
	GIN_MODE=release go build -tags netgo -ldflags '-s -w' -o app

clean:
	rm -f app

run: build
	GIN_MODE=release ./app
