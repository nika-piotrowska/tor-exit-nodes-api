run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run

redis:
	docker start tor-redis || docker run -d --name tor-redis -p 6380:6379 redis:7
