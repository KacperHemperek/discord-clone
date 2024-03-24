#!make
include .env

build:
	@go build -o ./bin/main

run: build
	@./bin/main

dev:
	air

test:
	go test -count=1 ./...