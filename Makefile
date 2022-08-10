rest: down up
up: run
down: stop
run:
	docker-compose -f ./docker-compose.yml up -d
stop:
	docker-compose -f ./docker-compose.yml down --remove-orphans
rebuild:
	docker-compose -f ./docker-compose.yml up -d --build
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
	go generate ./...
# for example `make cli command=addNetworkToWhitelist args="--n=192.168.0.0/24"`.
# For see all commands use `make cli`.
# Help for command `make cli command=help args=addNetworkToBlacklist`.
cli:
	docker-compose -f ./docker-compose.yml exec client ./abf-client $(command) $(args)
