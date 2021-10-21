MAIN_SERVICE_BINARY=main_service

PROJECT_DIR := ${CURDIR}

DOCKER_DIR := ${CURDIR}/docker

## build: Build compiles project
build:
	go build -o ${MAIN_SERVICE_BINARY} cmd/dripapp/main.go

## build-docker: Builds all docker containers
build-docker:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .

## run-and-build: Build and run docker
build-and-run: build-docker
	docker-compose up
	
## test-coverage: get final code coverage
test-coverage:
	go test -coverprofile=coverage.out.tmp ./...
	go tool cover -html=coverage.out.tmp

## test: launch all tests
test:
	go test ./...

## run-background: run process in background(available after build)
run-background:
	docker-compose up -d

## linter: linterint all files
linter: 
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused --disable deadcode
	go mod tidy

## run: Build and run docker with new changes
run: 
	docker rm -vf $$(docker ps -a -q) || true
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker-compose up --build --no-deps

## get: get all dependencies
get:
	go get ./...

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo