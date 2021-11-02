MAIN_SERVICE_BINARY=main_service

PROJECT_DIR := ${CURDIR}

DOCKER_DIR := ${CURDIR}/docker

## install-dependencies: Install docker
install-dependencies:
	sudo apt install docker.io

## build-go: Build compiles project
build-go:
	go mod tidy
	go build -o ${MAIN_SERVICE_BINARY} cmd/dripapp/main.go

## build-docker: Builds all docker containers
build-docker:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .

## test-coverage: get final code coverage
test-coverage:
	go test -coverprofile=coverage.out.tmp ./...
	go tool cover -html=coverage.out.tmp

## test: launch all tests
test:
	go test -coverprofile=coverage.out.tmp ./...

## run-background: run process in background(available after build)
run-background:
	docker-compose up --build --no-deps -d

## linter: linterint all files
linter:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused --disable deadcode
	go mod tidy

## build: Build and start docker with new changes
build:
	docker rm -vf $$(docker ps -a -q) || true
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .

## run: Run app
run:
# docker-compose up --build --no-deps
	docker-compose up --build --no-deps -d
	go run cmd/dripapp/main.go

## app: Build and run app
app:
# build run
	docker rm -vf $$(docker ps -a -q) || true
# rm -rf media
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker-compose up --build --no-deps -d
	go run cmd/dripapp/main.go

down:
	docker-compose down

## get: get all dependencies
get:
	go mod tidy
	go get ./...

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
