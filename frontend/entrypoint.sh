#!/bin/sh
set -e

echo "Generating runtime environment configuration..."

# Generate env-config.js with actual runtime values
cat > /usr/share/nginx/html/env-config.js <<EOF
window._env_ = {
  REACT_APP_API_URL: "${REACT_APP_API_URL:-http://localhost:8080}",
  REACT_APP_WS_URL: "${REACT_APP_WS_URL:-ws://localhost:8080}",
  REACT_APP_REFRESH_INTERVAL: "${REACT_APP_REFRESH_INTERVAL:-5000}",
  REACT_APP_VERSION: "${REACT_APP_VERSION:-1.0.0}",
  REACT_APP_LOGO_URL: "${REACT_APP_LOGO_URL:-}",
  REACT_APP_GITHUB_URL: "${REACT_APP_GITHUB_URL:-}",
  REACT_APP_LINKEDIN_URL: "${REACT_APP_LINKEDIN_URL:-}",
  REACT_APP_INFRA_OS: "${REACT_APP_INFRA_OS:-Talos}",
  REACT_APP_INFRA_STORAGE: "${REACT_APP_INFRA_STORAGE:-Longhorn/OpenEBS}",
  REACT_APP_INFRA_CNI: "${REACT_APP_INFRA_CNI:-Cilium}",
  REACT_APP_INFRA_LOAD_BALANCER: "${REACT_APP_INFRA_LOAD_BALANCER:-Traefik}",
  REACT_APP_INFRA_GITOPS: "${REACT_APP_INFRA_GITOPS:-Flux CD}"
};
EOF

echo "Runtime environment configuration generated:"
cat /usr/share/nginx/html/env-config.js

# Start nginx
exec nginx -g "daemon off;"