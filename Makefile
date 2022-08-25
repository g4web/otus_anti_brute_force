include ./configs/config.env

rest: down up
up: run
down: stop
run:
	docker-compose -f ./docker/docker-compose.yml --env-file ./configs/config.env up -d
stop:
	docker-compose -f ./docker/docker-compose.yml down --remove-orphans
build:
	docker-compose -f ./docker/docker-compose.yml up -d --build
test:
	go test -race ./...
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
	docker-compose -f ./docker/docker-compose.yml exec client ./abf-client $(command) $(args)

# for example `make migrate-create table=event`
migrate-create:
	migrate create -ext sql -dir ./internal/storage/sql/migrations $(table)
migrate-up:
	migrate -source file://internal/storage/sql/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/abf?sslmode=disable up
migrate-down:
	migrate -source file://internal/storage/sql/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/abf?sslmode=disable down 1
