BINARY_NAME=smshog

build:
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY_NAME}-darwin-apple-silicon main.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows main.go

clean:
	go clean
	rm -f bin/${BINARY_NAME}-linux
	rm -f bin/${BINARY_NAME}-darwin-apple-silicon
	rm -f bin/${BINARY_NAME}-windows