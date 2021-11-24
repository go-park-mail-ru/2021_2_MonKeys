MAIN_SERVICE_BINARY=main_service
CHAT_SERVICE_BINARY=chat_service
AUTH_SERVICE_BINARY=auth_service

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
	go build -o ${CHAT_SERVICE_BINARY} cmd/chat/main.go
	go build -o ${AUTH_SERVICE_BINARY} cmd/auth/main.go

## build-docker: Builds all docker containers
build-docker:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker build -t chat_service -f ${DOCKER_DIR}/chat_service.Dockerfile .
	docker build -t auth_service -f ${DOCKER_DIR}/chat_service.Dockerfile .

## test-coverage: get final code coverage
test-coverage:
	go test -coverprofile=coverage.out.tmp -coverpkg=./...  ./...
	cat coverage.out.tmp | grep -v mock > coverage2.out.tmp
	go tool cover -func=coverage2.out.tmp
	go tool cover -html=coverage2.out.tmp

test:
	go test ./...

## linter: linterint all files
linter:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused --disable deadcode
	go mod tidy

## clean: Clean all volumes, containers, media folders  and log files
clean:
	sudo service postgresql stop
	docker rm -vf $$(docker ps -a -q) || true
	sudo rm -rf media
	sudo rm -rf logs.log
	docker volume prune


##################################### deploy

## deploy-build: Deply build and start docker with new changes
deploy-build:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker build -t chat_service -f ${DOCKER_DIR}/chat_service.Dockerfile .
	docker build -t auth_service -f ${DOCKER_DIR}/chat_service.Dockerfile .

## deploy-run: Deploy run app
deploy-run:
	docker-compose -f prod.yml up --build --no-deps -d

## deploy-app: Deploy build and run app
deploy-app: deploy-build deploy-run

## deploy-app-clean: Deploy build and run app with clean
deploy-app-clean: clean deploy-build deploy-run

######################################## local

## build: Build and start docker with new changes
build:
	
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .

## run: Run app
run:
	docker-compose -f local.yml up --build --no-deps -d
	go run cmd/dripapp/main.go

run-chat:
	go run cmd/chat/main.go

run-auth:
	go run cmd/auth/main.go

## app: Build and run app
app: build run

## app-clean: Build and run app with clean
app-clean:clean build run


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
