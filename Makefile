MAIN_SERVICE_BINARY=main_service

PROJECT_DIR := ${CURDIR}

DOCKER_DIR := ${CURDIR}/docker

## install-dependencies: Install docker
install-dependencies:
	sudo apt install docker.io
	sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
	sudo chmod +x /usr/local/bin/docker-compose

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
	go test -coverprofile=coverage.out.tmp -coverpkg=./...  ./...
	cat coverage.out.tmp | grep -v mock > coverage2.out.tmp
	go tool cover -html=coverage2.out.tmp

test:
	go test ./...

## run-background: run process in background(available after build)
run-background:
	docker-compose up --build --no-deps -d

## linter: linterint all files
linter:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused --disable deadcode
	go mod tidy

## deploy-build: Deply build and start docker with new changes
deploy-build:
	docker rm -vf $$(docker ps -a -q) || true
	rm -rf media
	sudo rm -rf logs.log
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .

## deploy-run: Deploy run app
deploy-run:
	docker-compose -f prod.yml up --build --no-deps

## clean: Clean all volumes, containers, media folders  and log files
clean:
	docker rm -vf $$(docker ps -a -q) || true
	sudo rm -rf media
	sudo rm -rf logs.log
	docker volume prune

## build: Build and start docker with new changes
build:
	docker rm -vf $$(docker ps -a -q) || true
	sudo rm -rf media
	sudo rm -rf logs.log
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker-compose -f local.yml up --build --no-deps -d

## run: Run app
run:
	go run cmd/dripapp/main.go

## app: Build and run app
deploy-app: deploy-build deploy-run

app: build run


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
