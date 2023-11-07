BIN_MIGRATE := "./bin/migrate"

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			api/RateLimiter.proto

build-migrations:
	go build -v -o $(BIN_MIGRATE) ./cmd/migrations

migrate-status: build-migrations
	$(BIN_MIGRATE) -config="./configs/config.yml" status

migrate-up: build-migrations
	$(BIN_MIGRATE) -config="./configs/config.yml" up

migrate-down: build-migrations
	$(BIN_MIGRATE) -config="./configs/config.yml" down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...
