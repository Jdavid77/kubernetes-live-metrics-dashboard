# Serve with nginx (unprivileged)
FROM nginxinc/nginx-unprivileged:mainline-alpine-slim as unprivileged

# Update Alpine packages to fix security vulnerabilities
USER root
RUN apk update && \
    apk upgrade --no-cache && \
    rm -rf /var/cache/apk/*
USER nginx

# Build
FROM node:26-alpine AS builder

WORKDIR /app

COPY package.json ./
RUN npm install
COPY . .

RUN npm run build


FROM unprivileged

# Use the existing nginx user (uid/gid 101) by name
COPY --chown=nginx:nginx --from=builder /app/build /usr/share/nginx/html
# The base image's entrypoint renders templates with envsubst at startup:
# /etc/nginx/templates/nginx.conf.template -> /etc/nginx/nginx.conf
COPY --chown=nginx:nginx nginx.conf /etc/nginx/templates/nginx.conf.template
ENV NGINX_ENVSUBST_OUTPUT_DIR=/etc/nginx \
    NGINX_ENVSUBST_FILTER=^NGINX_BACKEND_URL$ \
    NGINX_BACKEND_URL=localhost:8000

EXPOSE 8080