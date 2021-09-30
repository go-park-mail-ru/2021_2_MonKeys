clean: 
	rm -rf build
	
test-coverage:
	go test -coverprofile=coverage.out.tmp ./...
	go tool cover -html=coverage.out.tmp

test: 
	go test ./...

run: 
	go run server

build:
	go get ./...
	go build -v -o build/apiServer server