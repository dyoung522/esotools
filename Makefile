fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet build-osx build-linux build-win

build-osx:
	GOOS=darwin GOARCH=arm64 go build -o "bin/$(shell basename $(PWD))-darwin-arm64" ./main.go
	GOOS=darwin GOARCH=amd64 go build -o "bin/$(shell basename $(PWD))-darwin-amd64" ./main.go

build-win:
	GOOS=windows GOARCH=amd64 go build -o "bin/$(shell basename $(PWD)).exe" ./main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o  "bin/$(shell basename $(PWD))-linux" ./main.go
