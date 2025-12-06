# Base
FROM node:22-alpine AS base

RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

# Build
FROM base AS builder

WORKDIR /app

COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

# Serve with nginx (unprivileged)
FROM nginxinc/nginx-unprivileged:stable-alpine as nginx

# Use the existing nginx user (uid/gid 101) by name
COPY --chown=nginx:nginx --from=builder /app/build /usr/share/nginx/html
COPY --chown=nginx:nginx nginx.conf /etc/nginx/nginx.conf

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]