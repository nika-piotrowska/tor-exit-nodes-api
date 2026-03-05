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
(coming soon)

## Configuration (planned)
- TOR_EXIT_ADDRESSES_URL – URL to Tor exit addresses list
- REDIS_URL – Redis address (default: redis://localhost:6380/0)
- TOR_EXIT_CACHE_TTL – cache time-to-live in seconds for the Tor exit nodes list (default: 300)