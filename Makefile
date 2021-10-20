clean: 
	rm -rf build
	
test-coverage:
	go test -coverprofile=coverage.out.tmp ./...
	go tool cover -html=coverage.out.tmp

test: 
	go test ./...

linter: 
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused --disable deadcode
	go mod tidy

run: 
	go run cmd/dripapp/main.go
	# build/dripapp

get:
	go get ./...

build:
	rm -rf build/dripapp
	go build -v -o ./build/dripapp cmd/dripapp/main.go