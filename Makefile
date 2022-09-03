include ./configs/config.env
SERVER := "server"
CLIENT := "client"

rest: down up
up: build run
down: stop
run:
	docker-compose -f ./docker/docker-compose.yml --env-file ./configs/config.env up -d
stop:
	docker-compose -f ./docker/docker-compose.yml --env-file ./configs/config.env down --remove-orphans
build-docker:
	docker-compose -f ./docker/docker-compose.yml --env-file ./configs/config.env up -d --build
build:
	go build -v -o ./bin/$(SERVER) ./cmd/$(SERVER)
	go build -v -o ./bin/$(CLIENT) ./cmd/$(CLIENT)

unit-tests:
	go test -race ./internal/...

integration-tests: build
	docker-compose -f ./docker/docker-compose.integr.yml --env-file ./configs/config_test.env up -d --build;\
	test_status_code=0 ;\
	docker-compose -f ./docker/docker-compose.integr.yml --env-file ./configs/config_test.env run integr go test -race /abf/test... || test_status_code=$$? ;\
	docker-compose -f ./docker/docker-compose.integr.yml --env-file ./configs/config_test.env down --remove-orphans;\
	exit $$test_status_code ;

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...
fmt:
	gofmt -s -w .
	gofumpt -l -w .
generate:
	protoc ./api/abf.proto --go_out=./internal/proto/ --go-grpc_out=./internal/proto/
# for example `make cli command=addNetworkToWhitelist args="--n=192.168.0.0/24"`.
# For see all commands use `make cli`.
# Help for command `make cli command=help args=addNetworkToBlacklist`.
cli:
	docker-compose -f ./docker/docker-compose.yml exec client /abf/bin/client $(command) $(args)

# for example `make migrate-create table=event`
migrate-create-old:
	migrate create -ext sql -dir ./internal/storage/sql/migrations $(table)
migrate-up-old:
	migrate -source file://internal/storage/sql/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/abf?sslmode=disable up
migrate-down-old:
	migrate -source file://internal/storage/sql/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/abf?sslmode=disable down 1


# for example `make migrate-create table=event`
migrate-create:
	docker-compose -f ./docker/docker-compose.yml  --env-file ./configs/config.env run migration /go/bin/goose postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"  create $(table) sql
migrate-up:
	docker-compose -f ./docker/docker-compose.yml  --env-file ./configs/config.env run migration /go/bin/goose postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"  up
migrate-down:
	docker-compose -f ./docker/docker-compose.yml  --env-file ./configs/config.env run migration /go/bin/goose postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"  down
