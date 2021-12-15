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
	go build -o ${MAIN_SERVICE_BINARY} cmd/dripapp/main.go
	go build -o ${CHAT_SERVICE_BINARY} cmd/chat/main.go
	go build -o ${AUTH_SERVICE_BINARY} cmd/auth/main.go

## test-coverage: get final code coverage
test-coverage:
	go test -coverprofile=coverage.out.tmp -coverpkg=./...  ./...
	cat coverage.out.tmp | grep -v mock | grep -v cmd | grep -v easyjson > coverage2.out.tmp
	go tool cover -func=coverage2.out.tmp
	go tool cover -html=coverage2.out.tmp -o cover.html

## test: test code
test:
	go test ./...

## linter: linterint all files
linter:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --disable unused
	go mod tidy

## clean: Clean all volumes, containers, media folders  and log files
clean:
	sudo service postgresql stop
	docker rm -vf $$(docker ps -a -q) || true
	sudo rm -rf media
	sudo rm -rf logs.log
	docker volume prune

clean-deploy:
	sudo docker-compose rm postgres
	sudo docker rm -vf $$(docker ps -a -q) || true
	sudo rm -rf media
	sudo rm -rf logs.log
	sudo docker volume prune -f

##################################### deploy

## deploy-build: Deply build and start docker with new changes
build:
	cat prod.json > config.json
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker build -t chat_service -f ${DOCKER_DIR}/chat_service.Dockerfile .
	docker build -t auth_service -f ${DOCKER_DIR}/auth_service.Dockerfile .


## deploy-run: Deploy run app on background
deploy-run:
	sudo docker-compose -f prod.yml up --build --no-deps -d

## deploy: Deploy build and run app
deploy: build deploy-run

## redeploy: Deploy build and run app with clean
redeploy: clean-deploy build deploy-run

######################################## local

local-build:
	cat local.json > config.json
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker build -t chat_service -f ${DOCKER_DIR}/chat_service.Dockerfile .
	docker build -t auth_service -f ${DOCKER_DIR}/auth_service.Dockerfile .

## deploy-run: Deploy run app on background
local-run:
	docker-compose -f local.yml up --build --no-deps

local: local-build local-run

relocal: clean local-build local-run

######################################## debug

## build: Build and start docker with new changes
debug:
	cat debug.json > config.json
	docker build -t drip_tarantool -f ${DOCKER_DIR}/drip_tarantool.Dockerfile .
	docker-compose -f debug.yml up --build --no-deps -d

run-dripapp:
	go run cmd/dripapp/main.go

run-chat:
	go run cmd/chat/main.go

run-auth:
	go run cmd/auth/main.go

## redebug: Build and debug app with clean
redebug:clean debug

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
