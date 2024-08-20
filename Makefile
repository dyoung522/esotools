tidy:
	go mod tidy

fmt: tidy
	go fmt ./...

vet: fmt
	go vet ./...

build: vet build-osx build-linux build-win

build-osx:
	GOOS=darwin GOARCH=amd64 go build -o "bin/$(shell basename $(PWD))-darwin-amd64" ./main.go

build-win:
	GOOS=windows GOARCH=amd64 go build -o "bin/$(shell basename $(PWD)).exe" ./main.go
