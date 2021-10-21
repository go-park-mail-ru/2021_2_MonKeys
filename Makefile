MAIN_SERVICE_BINARY=main_service

PROJECT_DIR := ${CURDIR}

DOCKER_DIR := ${CURDIR}/docker

## build: Build compiles project
build:
	go build -o ${MAIN_SERVICE_BINARY} cmd/dripapp/main.go

build-docker:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .

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