#!make
include .env

build:
	@go build -o ./bin/main

run: build
	@./bin/main

test:
	go test -count=1 ./...