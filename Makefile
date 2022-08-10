rest: down up
up:
	docker-compose -f ./docker-compose.yml up -d
down:
	docker-compose -f ./docker-compose.yml down --remove-orphans

test:
	go test -race ./...
	#docker-compose -f ./docker-compose.yml exec abf2 go test -race ./...

#lint:
#	go install github.com/golangci/golangci-lint/cmd/golangci-lint
#	golangci-lint run ./...

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	gofumpt -l -w .