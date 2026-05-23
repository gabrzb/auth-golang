FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM caddy:2-alpine

RUN addgroup -S caddy && adduser -S -D -H -h /srv -s /sbin/nologin -G caddy caddy \
  && chown -R caddy:caddy /data /config /usr/share/caddy /etc/caddy

COPY --chown=caddy:caddy Caddyfile /etc/caddy/Caddyfile
COPY --chown=caddy:caddy --from=builder /app/dist /usr/share/caddy

EXPOSE 80

USER caddy
