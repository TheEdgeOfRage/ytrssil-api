.PHONY: all setup ytrssil-api build gen-mocks lint yamllint test migrate image-build image-push

DB_URI ?= postgres://ytrssil:ytrssil@localhost:5431/ytrssil?sslmode=disable

all: lint test build

setup: bin/golangci-lint
	go mod download

ytrssil-api:
	go build -o dist/ytrssil-api cmd/main.go

build: ytrssil-api

bin/moq:
	GOBIN=$(PWD)/bin go install github.com/matryer/moq@v0.2.7

gen-mocks: bin/moq
	./bin/moq -pkg db_mock -out ./mocks/db/db.go ./db DB
	go fmt ./...

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.48.0

lint: bin/golangci-lint
	go fmt ./...
	go vet ./...
	bin/golangci-lint -c .golangci.yml run ./...
	go mod tidy

test:
	go mod tidy
	go test -timeout=10s -race -benchmem ./...

migrate:
	migrate -database "$(DB_URI)" -path migrations up

image-build:
	@echo "# Building docker image..."
	docker build -t theedgeofrage/ytrssil:api --target api .
	docker build -t theedgeofrage/ytrssil:migrations --target migrations .

image-push: image-build
	@echo "# Pushing docker image..."
	docker push theedgeofrage/ytrssil:api
	docker push theedgeofrage/ytrssil:migrations
