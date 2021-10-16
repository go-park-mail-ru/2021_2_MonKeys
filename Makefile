clean: 
	rm -rf build
	
test-coverage:
	go test -coverprofile=coverage.out.tmp ./...
	go tool cover -html=coverage.out.tmp

test: 
	go test ./...

linter: 
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused --disable deadcode

run: 
	go run server

get:
	go get ./...

get-linter:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1

build:
	go build -v -o ./build/apiDripapp cmd/dripapp/main.go