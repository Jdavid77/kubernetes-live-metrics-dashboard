# Kubernetes Dashboard

A modern, real-time Kubernetes cluster monitoring dashboard built with Go and React.

## Features

- **Real-time Metrics**: Live cluster metrics updated every 5 seconds via WebSocket
  - CPU Usage
  - Memory Usage
  - Running Pods Count
  - Ingresses Count

- **Application Monitoring**: View all deployed applications with:
  - Health status (Healthy/Degraded/Down)
  - Pod counts
  - Technology stack detection
  - Namespace information

- **Modern UI**: Dark-themed interface with gradient backgrounds and smooth animations

- **Docker Support**: Fully containerized with Docker Compose

## Architecture

### Backend (Go)
- **Auto-detection**: Supports both in-cluster and local kubeconfig connections
- **Metrics Collection**: Integrates with Kubernetes metrics-server
- **WebSocket**: Real-time updates with automatic reconnection
- **Caching**: 5-second cache to reduce API calls
- **CORS**: Configurable CORS middleware

### Frontend (React)
- **React Hooks**: Custom hooks for data management
- **WebSocket Client**: Real-time data with polling fallback
- **Responsive Design**: Mobile-friendly interface
- **Error Handling**: Comprehensive error boundaries

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- Docker and Docker Compose (optional)
- Access to a Kubernetes cluster
- Kubernetes metrics-server installed in the cluster

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/Jdavid77/kubernetes-dashboard.git
cd kubernetes-dashboard
```

2. Ensure your kubeconfig is at `~/.kube/config`

3. Start the application:
```bash
docker-compose up --build
```

4. Access the dashboard:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Local Development

#### Backend

```bash
cd backend

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm start
```

## Configuration

### Backend Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `KUBECONFIG` | `""` | Path to kubeconfig (auto-detected if empty) |
| `CORS_ORIGIN` | `http://localhost:3000` | Allowed CORS origin |
| `METRICS_REFRESH_INTERVAL` | `5s` | Metrics refresh interval |
| `LOG_LEVEL` | `info` | Logging level |

### Frontend Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `REACT_APP_API_URL` | `http://localhost:8080` | Backend API URL |
| `REACT_APP_WS_URL` | `ws://localhost:8080` | WebSocket URL |
| `REACT_APP_REFRESH_INTERVAL` | `5000` | Fallback polling interval (ms) |

## Deploying to Kubernetes

### 1. Create RBAC Resources

```bash
kubectl apply -f rbac.yaml
```

### 2. Deploy the Application

See the included `rbac.yaml` for the required ServiceAccount, ClusterRole, and ClusterRoleBinding.

The application requires read-only access to:
- Pods, Nodes (core API)
- Deployments (apps API)
- Ingresses (networking.k8s.io API)
- Pod and Node metrics (metrics.k8s.io API)

### 3. Access the Dashboard

Port-forward or use an Ingress to access the dashboard externally.

## API Endpoints

### REST Endpoints

- `GET /api/health` - Health check
- `GET /api/metrics` - Get current cluster metrics
- `GET /api/applications` - List all applications
- `GET /api/applications/:name?namespace=<ns>` - Get application details

### WebSocket Endpoint

- `GET /ws/metrics` - Real-time metrics stream

## Project Structure

```
kubernetes-dashboard/
├── backend/              # Go backend
│   ├── cmd/api/         # Application entry point
│   ├── internal/        # Internal packages
│   │   ├── config/     # Configuration
│   │   ├── handlers/   # HTTP handlers
│   │   ├── kubernetes/ # K8s client & operations
│   │   ├── middleware/ # HTTP middleware
│   │   ├── models/     # Data models
│   │   └── server/     # Server setup
│   └── Dockerfile
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/ # React components
│   │   ├── hooks/      # Custom hooks
│   │   ├── services/   # API services
│   │   ├── styles/     # Global styles
│   │   └── utils/      # Utilities
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
└── README.md
```

## Development

### Backend

```bash
cd backend

# Run tests
go test ./...

# Build binary
go build -o api ./cmd/api

# Run with live reload (install air first)
air
```

### Frontend

```bash
cd frontend

# Run tests
npm test

# Build for production
npm run build

# Lint code
npm run lint
```

## Troubleshooting

### WebSocket Connection Failed

If WebSocket connection fails, the dashboard will automatically fall back to HTTP polling. Check:
- CORS configuration
- Firewall rules
- WebSocket support in your environment

### No Metrics Available

If metrics show as "0%" or "0GB":
- Ensure metrics-server is installed: `kubectl get deployment metrics-server -n kube-system`
- Check metrics API: `kubectl top nodes`

### Permission Denied Errors

Ensure the ServiceAccount has proper RBAC permissions. See `rbac.yaml` for required permissions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - See LICENSE file for details

## Acknowledgments

Inspired by [catdevops.net](https://catdevops.net)
