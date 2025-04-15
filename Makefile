-include .env
export

APP=whereismycity
CURRENT_DIR=$(shell pwd)
CMD_DIR=./cmd

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/main.go

# migrate
.PHONY: migrate
migrate:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

# generate swagger
.PHONY: swagger-gen
swagger-gen:
	swag init --parseDependency --dir ./delivery/api -g router.go -o ./delivery/api/docs