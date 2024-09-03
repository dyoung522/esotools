default: check

## Main Commands

build: fmt clean test build-osx build-linux build-win
clean: clean-bin tidy

## Supporting Commands

tidy:
	go mod tidy

fmt: tidy
	trunk fmt

fmt-all: tidy
	trunk fmt --all

check: fmt
	trunk check

check-all: fmt-all
	trunk check --all

test:
	ginkgo -r

clean-bin:
	rm -f bin/*

## Build sub-commands

build-osx:
	GOOS=darwin GOARCH=amd64 go build -o "bin/$(shell basename $(PWD))-osx" ./main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o "bin/$(shell basename $(PWD))-linux" ./main.go

build-win:
	GOOS=windows GOARCH=amd64 go build -o "bin/$(shell basename $(PWD)).exe" ./main.go

## Git Hooks
pre-commit: clean check test
	git add go.mod go.sum
