# Tor Exit Nodes API (Go)

A small HTTP API service that fetches the [Tor](https://www.torproject.org/) [exit node list](https://check.torproject.org/exit-addresses), caches it in Redis, and exposes endpoints to check whether a given IP address is a Tor exit node.

## Tech stack
- Go
- Redis (via Docker)
- HTTP API

## Endpoints (planned)
- GET /up – health check
- GET /api/v1/{ip} – returns 200 and JSON with the IP if it is a Tor exit node, 404 otherwise
- GET /api/status – returns cache status (exists, size, refreshed_at)

## Local development

### Run Redis
```bash
docker run -d --name tor-redis -p 6380:6379 redis:7
```

### Run the server
```bash
make run
```

## Linters (golangci-lint)

This repository uses `golangci-lint` to keep code quality consistent.

### Step-by-step setup

1. Install golangci-lint:
```bash
make lint-install
```

2. Verify it is installed:
```bash
golangci-lint version
```

3. Run linters:
```bash
make lint
```

4. Run linters with automatic fixes where possible:
```bash
make lint-fix
```

### Useful commands
- `make lint` → run all configured linters
- `make lint-fix` → apply auto-fixes and re-run checks
- `golangci-lint run ./...` → direct linter invocation

The linter configuration lives in `.golangci.yml`.

## Releases

This repository uses semantic versioning and Git tags for releases.

To create a new release:

```bash
./scripts/bump-version.sh X.Y.Z --push
```

This command will:

- update the `VERSION` file
- prepend a new entry to `CHANGELOG.md`
- create a commit `release: vX.Y.Z`
- create an annotated Git tag (`vX.Y.Z`)
- push the commit and tag to GitHub

Example:

```bash
./scripts/bump-version.sh 0.2.0 --push
```

## Configuration (planned)
- TOR_EXIT_ADDRESSES_URL – URL to Tor exit addresses list
- REDIS_URL – Redis address (default: redis://localhost:6380/0)
- TOR_EXIT_CACHE_TTL – cache time-to-live in seconds for the Tor exit nodes list (default: 300)
