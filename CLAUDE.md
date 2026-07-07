# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This App Does

Real-time Kubernetes cluster monitoring dashboard. A Go backend connects to the Kubernetes API, aggregates live metrics (CPU, memory, pods, nodes, deployments), and pushes updates over WebSocket. A React SPA visualises this data.

## Commands

### Using Make (recommended)

```bash
make install          # install all deps (backend go mod + frontend npm)
make run              # run backend + frontend dev servers in parallel
make build            # build both for production
make test             # run all tests (backend + frontend)

make backend-build    # go build -o api ./cmd/api
make backend-test     # go test -v ./...
make backend-run      # go run cmd/api/main.go

make frontend-install # npm install
make frontend-dev     # npm start (React dev server on :3000)
make frontend-build   # npm run build
make frontend-test    # npm test
```

### Docker (recommended for full stack)

```bash
make docker-build     # docker-compose build
make docker-up        # docker-compose up -d
make docker-down      # docker-compose down
make docker-logs      # docker-compose logs -f
```

Requires `~/.kube/config` pointing at a cluster with `metrics-server` installed. Frontend served at `http://localhost:3000`, backend API at `http://localhost:8000`.

### Kubernetes RBAC

```bash
make k8s-apply-rbac   # kubectl apply -f rbac.yaml
make k8s-delete-rbac  # kubectl delete -f rbac.yaml
```

### Single backend test

```bash
cd backend && go test -v ./internal/kubernetes/...
cd backend && go test -run TestFunctionName ./...
```

## Architecture

### Data Flow

```
Browser (React SPA on :3000)
  ├── WebSocket  /ws/metrics   → Go backend (push every 5s)
  └── HTTP REST  /api/*        → Go backend

Go backend (:8080)
  └── client-go + k8s/metrics → Kubernetes API server
```

In Docker/production, nginx reverse-proxies `/api/` and `/ws/` to `http://localhost:8000` (backend), so the browser talks only to nginx. There is no CORS issue at runtime.

### Backend (`/backend`)

Entry point: `cmd/api/main.go` — wires config → k8s client → aggregator → handlers → gorilla/mux router.

Key packages:
- `internal/config` — env vars: `PORT` (8080), `KUBECONFIG`, `CORS_ORIGIN`, `METRICS_REFRESH_INTERVAL` (5s)
- `internal/kubernetes/client.go` — auto-detects in-cluster config, falls back to `~/.kube/config`; two clientsets (standard + metrics)
- `internal/kubernetes/source.go` — `MetricsSource` interface; `*Aggregator` is the prod adapter; implement a fake to test handlers without a cluster
- `internal/kubernetes/aggregator.go` — in-memory cache with 5s TTL and `sync.RWMutex`; all six public methods (`GetClusterMetrics`, `GetApplications`, `GetApplicationDetail`, `GetPods`, `GetServices`, `GetNodes`, `GetNamespaces`) are cached
- `internal/handlers/websocket.go` — `Hub` struct (register/broadcast/shutdown); started with a `context.Context` so goroutines stop on SIGINT/SIGTERM; `WebSocketHandler` is a thin HTTP→WebSocket upgrade adapter
- `internal/server/router.go` — gorilla/mux with Logger + CORS middleware
- `internal/server/server.go` — graceful shutdown on SIGINT/SIGTERM (30s timeout)
- `internal/models/` — plain Go structs: `ClusterMetrics`, `Application`, `Pod`, `Service`, `Node`

There is no database — all state is live-fetched from the Kubernetes API; only in-memory caching.

### Frontend (`/frontend`)

CRA (Create React App) + React 18. Key files:
- `src/App.jsx` — root: `Header` + `Dashboard` inside `ErrorBoundary`
- `src/components/Dashboard/` — `ClusterOverview` (metric tiles), `NodesSection` (node cards), `Stack` (infra labels from env)
- `src/components/ApplicationGrid/` — application cards with health status
- `src/hooks/useMetrics.js` — subscribes to `/ws/metrics` via `useWebSocket`; falls back to HTTP polling
- `src/hooks/useWebSocket.js` — WebSocket lifecycle with exponential-backoff reconnect (max 5 attempts)
- `src/services/` — Axios REST clients for metrics, applications, pods, services, nodes
- `src/config/env.js` — merges `window._env_` (Docker/production) and `process.env` (dev)

### Runtime Configuration

Frontend env vars are injected at container start by `frontend/entrypoint.sh`, which writes `frontend/public/env-config.js` from container env vars as `window._env_`. In dev mode, use `REACT_APP_*` prefixes in `.env`.

### CI / Release

Conventional commits on `main` trigger `semantic-release` (`.releaserc.js`), which bumps `frontend/package.json` version, generates CHANGELOG, creates a GitHub release, then triggers a multi-arch Docker image build pushed to GHCR (`ghcr.io/Jdavid77/kubernetes-live-metrics-dashboard-{backend,frontend}`).
