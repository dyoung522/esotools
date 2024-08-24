tidy:
	go mod tidy

fmt: tidy
	go fmt ./...

vet: fmt
	${GOBIN}/staticcheck ./...

clean_bin:
	rm -f bin/*

clean: clean_bin tidy

build: vet clean_bin build-osx build-win

build-osx:
	GOOS=darwin GOARCH=amd64 go build -o "bin/$(shell basename $(PWD))" ./main.go

build-win:
	GOOS=windows GOARCH=amd64 go build -o "bin/$(shell basename $(PWD)).exe" ./main.go

test:
	ginkgo -r
