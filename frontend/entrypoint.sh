#!/bin/sh
set -e

echo "Generating runtime environment configuration..."

# Read version from version.txt file (created during build from package.json)
VERSION=$(cat /version.txt 2>/dev/null || echo "1.0.0")

# Generate env-config.js with actual runtime values
cat > /usr/share/nginx/html/env-config.js <<EOF
window._env_ = {
  REACT_APP_REFRESH_INTERVAL: "${REACT_APP_REFRESH_INTERVAL:-5000}",
  REACT_APP_VERSION: "${VERSION}",
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

# Start nginx
exec nginx -g "daemon off;"