APP := tor-exit-nodes-api # := natychmiastowe przypisanie wartości (w momencie definicji, wykonuje się raz - wartość się nie zmieni w trakcie działania make)
CMD := ./cmd/api # DRY
BIN := bin/$(APP)

GOLANGCI_LINT_VERSION := v1.57.2 # uspójnienie wersji linterów crosszespołowo

# target nie reprezentuje pliku, tylko komendę, która powinna być zawsze wykonana. Make zawsze wykonuje ten target i nie sprawdza czy istnieje plik o tej samej nazwie
.PHONY: run build test lint lint-fix fmt fmt-check check lint-install redis redis-stop

run:
	go run $(CMD)

build:
	mkdir -p bin
	go build -o $(BIN) $(CMD)

test:
	go test ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

redis:
	docker start tor-redis || docker run -d --name tor-redis -p 6380:6379 redis:7

redis-stop:
	docker stop tor-redis