```md
# Shagram — Go WebSocket Chat | DevOps Pet Project

Shagram is a small real-time chat application written in Go.  
This is primarily a **DevOps-focused pet project**: the app itself is intentionally simple, while the main goal is to practice CI/CD, containerization, reverse proxy/TLS, and operating self-hosted infrastructure.

> Public repo note: I intentionally do not publish any real server URLs, IP addresses, credentials, or registry endpoints here.

## What it does
- Multi-room chat via WebSockets (`/ws/:room`) with message broadcast.
- Message history persisted in SQLite and available through an HTTP API.
- Minimal browser UI served from `./static`.

## DevOps skills demonstrated
- Cloud: Deployed and operated on a VPS in **Yandex Cloud** (self-managed infrastructure).
- Docker: Multi-stage image build for a Go service; small runtime image.
- Docker Compose: Stack orchestration for the app (Go service + Nginx) and for CI infrastructure (Jenkins controller + agent).
- Nginx: Reverse proxy configuration for HTTP + WebSocket (Upgrade/Connection headers) and TLS termination.
- TLS: Self-signed certificates for dev/demo environments; production note to use a trusted CA (e.g., Let’s Encrypt).
- CI/CD with Jenkins: Pipeline that builds a Docker image, tags it, pushes it to a registry, and deploys via Docker Compose (deploy gated to `main`).
- GitHub → Jenkins automation: GitHub repository **webhook** triggers Jenkins builds on push/changes.
- Private registry (Harbor): Self-hosted registry for storing and distributing built images.
- Jenkins agent architecture: Dedicated inbound Docker agent for builds, with Docker socket mounting for Docker-based workloads (security caveat applies).
- Ops automation: One-command restart/update script (`restart-stacks.sh`) that pulls and recreates Shagram, Jenkins, and Harbor stacks.

## Roadmap
- Kubernetes: Migrate deployment from Docker Compose to Kubernetes (manifests/Helm), add Ingress + cert-manager, and prepare the app for future scaling (multi-replica WebSocket strategy and persistent storage).

## Application stack
- Go + Gin (HTTP API and routing) (`cmd/server`).
- WebSockets: Gorilla WebSocket (`internal/api`, `internal/websocket`).
- Auth: JWT access tokens (`internal/auth`), login endpoint issues tokens, and WebSocket uses `?token=...`.
- Storage: SQLite (`internal/db`) initialized from `migrations/schema.sql`.

## Configuration
Environment variables:
- `JWT_SECRET` (required): signing key for JWT tokens.
- `DATABASE_PATH` (optional): SQLite file path, defaults to `/app/data/shagram.db` in the container.
- `WS_ALLOWED_ORIGINS` (required for WebSocket): comma-separated list of allowed `Origin` values for browser WebSocket connections.
- `APP_IMAGE` (optional, deployment): Docker image reference used by Compose to deploy a prebuilt image.

## HTTP API
Server listens on `:8080`.

- `POST /api/auth/login` → `{ "access_token": "..." }`
- `GET /api/me` (requires `Authorization: Bearer <token>`) → `{ "username": "..." }`
- `GET /api/rooms` → `{ "rooms": [...] }`
- `GET /api/messages/:room` → `{ "messages": [...] }` (last 50 messages)

## WebSocket
Endpoint:
- `ws(s)://<host>/ws/<room>?token=<access_token>`

Client sends JSON:
```json
{ "text": "hello" }
```

Server broadcasts plain text messages:
- `alice: hello`

Security notes:
- WebSocket requires a JWT token (`?token=...`) and validates the `Origin` header against `WS_ALLOWED_ORIGINS`.
- This project is still a learning lab; for production you would additionally harden auth, rate limits, and request validation.

## Data model (SQLite)
Schema is created on startup from `migrations/schema.sql` and includes:
- `rooms(id, name)`
- `messages(id, room_id, user, text, created_at)`

## Quickstart (local, without Docker)
Prerequisites: Go toolchain.

```bash
export JWT_SECRET=change-me
export DATABASE_PATH=./shagram.db
export WS_ALLOWED_ORIGINS=http://localhost:8080

go run ./cmd/server
```

Open:
- http://localhost:8080

## Quickstart (Docker Compose + Nginx TLS)
Prerequisites: Docker + Docker Compose.

1) Generate a self-signed TLS certificate (dev/demo only):
See `deploy/shagram/nginx/certs/README.md`.

2) Create `deploy/shagram/.env`:
```bash
JWT_SECRET=change-me
WS_ALLOWED_ORIGINS=https://localhost
# Optional: override the app image (e.g., from your registry)
# APP_IMAGE=<your-registry>/<project>/shagram:<tag>
```

3) Start:
```bash
cd deploy/shagram
docker compose up -d --build
```

Open:
- https://localhost

## CI/CD overview (Jenkins + Registry)
- Jenkins pipeline builds a Docker image from this repository, tags it, pushes it to a registry, and deploys the updated stack with Docker Compose (deploy gated to `main`).
- Jenkins is triggered by a GitHub repository webhook (push/changes).
- Jenkins setup instructions: `infra/jenkins/README.md`.

## Ops: restart / update all stacks
On the host, `restart-stacks.sh` can be used to pull and recreate:
- Shagram stack (app + nginx)
- Jenkins stack (controller + agent)
- Harbor stack (registry)

## Repository structure
- `cmd/server/` — Gin router, endpoints, wiring.
- `internal/auth/` — JWT issuing/parsing + Gin auth middleware.
- `internal/api/` — WebSocket handler (token validation + origin check + DB persistence + broadcast).
- `internal/websocket/` — In-memory hub/rooms and broadcast logic.
- `internal/db/` — SQLite connection + schema bootstrap.
- `migrations/` — SQL schema.
- `static/` — Minimal web UI.
- `deploy/shagram/` — Docker Compose + Nginx config/certs.
- `infra/jenkins/` — Jenkins controller/agent Compose setup and docs.

## Notes / trade-offs
- The WebSocket hub is in-memory, so the app is intended to run as a single instance (no horizontal scaling).
- Jenkins Docker agent mounts `/var/run/docker.sock`, which provides high-level control of the Docker host; use only in trusted environments.
- Self-signed TLS is for development/demo; production should use a trusted CA (e.g., Let’s Encrypt).
```
