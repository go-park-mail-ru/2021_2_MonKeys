clean: 
	rm -rf build
	
test-coverage:
	go test -coverprofile=coverage.out.tmp ./...
	go tool cover -html=coverage.out.tmp

test: 
	go test ./...

linter: 
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

run: 
	go run server

build:
	go get ./...
	go build -v -o build/apiServer server