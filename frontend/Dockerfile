# Serve with nginx (unprivileged)
FROM nginxinc/nginx-unprivileged:stable-alpine as unprivileged

# Update Alpine packages to fix security vulnerabilities
USER root
RUN apk update && \
    apk upgrade --no-cache && \
    rm -rf /var/cache/apk/*
USER nginx

# Build
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json ./
RUN npm install
COPY . .

RUN node -p "require('./package.json').version" > /tmp/version.txt
RUN REACT_APP_VERSION=$(cat /tmp/version.txt) npm run build


FROM unprivileged

# Use the existing nginx user (uid/gid 101) by name
COPY --chown=nginx:nginx --from=builder /app/build /usr/share/nginx/html
COPY --chown=nginx:nginx --from=builder /tmp/version.txt /version.txt
COPY --chown=nginx:nginx nginx.conf /etc/nginx/nginx.conf
COPY --chown=nginx:nginx entrypoint.sh /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]