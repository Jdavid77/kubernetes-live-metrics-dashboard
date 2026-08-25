#!/bin/sh
set -e

echo "Generating runtime environment configuration..."

# Generate env-config.js with actual runtime values
cat > /usr/share/nginx/html/env-config.js <<EOF
window._env_ = {
  REACT_APP_REFRESH_INTERVAL: "${REACT_APP_REFRESH_INTERVAL:-5000}",
  REACT_APP_LOGO_URL: "${REACT_APP_LOGO_URL:-}",
  REACT_APP_GITHUB_URL: "${REACT_APP_GITHUB_URL:-}",
  REACT_APP_LINKEDIN_URL: "${REACT_APP_LINKEDIN_URL:-}",
  REACT_APP_INFRA_OS: "${REACT_APP_INFRA_OS:-}",
  REACT_APP_INFRA_STORAGE: "${REACT_APP_INFRA_STORAGE:-}",
  REACT_APP_INFRA_CNI: "${REACT_APP_INFRA_CNI:-}",
  REACT_APP_INFRA_LOAD_BALANCER: "${REACT_APP_INFRA_LOAD_BALANCER:-}",
  REACT_APP_INFRA_GITOPS: "${REACT_APP_INFRA_GITOPS:-}"
};
EOF

echo "Runtime environment configuration generated:"
cat /usr/share/nginx/html/env-config.js

# Substitute NGINX_BACKEND_URL in nginx.conf (defaults to localhost:8000 for k8s sidecar)
NGINX_BACKEND_URL="${NGINX_BACKEND_URL:-localhost:8000}" \
  envsubst '${NGINX_BACKEND_URL}' < /etc/nginx/nginx.conf > /tmp/nginx.conf
cp /tmp/nginx.conf /etc/nginx/nginx.conf

# Start nginx
exec nginx -g "daemon off;"