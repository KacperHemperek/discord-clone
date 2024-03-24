#!make
include .env

build:
	@go build -o ./tmp/main

run: build
	@./tmp/main

dev:
	air

test:
	go test -count=1 ./...