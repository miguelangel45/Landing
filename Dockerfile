FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY config-server/go.mod .
RUN go mod download
COPY config-server/main.go .
RUN CGO_ENABLED=0 go build -o config-server .

FROM nginx:alpine

# Copy the landing page and assets
COPY index.html /usr/share/nginx/html/index.html
COPY *.png *.svg *.webp /usr/share/nginx/html/

# Copy the config server binary
COPY --from=builder /app/config-server /usr/local/bin/config-server

# Configure nginx: landing page + reverse proxy /api/config to the Go server
RUN printf 'server {\n\
    listen 8000;\n\
    root /usr/share/nginx/html;\n\
    index index.html;\n\
    gzip on;\n\
    gzip_types text/html text/css application/javascript image/svg+xml;\n\
    location /api/ {\n\
        proxy_pass http://127.0.0.1:8080;\n\
        proxy_set_header Host $host;\n\
        proxy_set_header X-Real-IP $remote_addr;\n\
    }\n\
    location ~* \\.(webp|png|svg|jpg|jpeg|ico)$ {\n\
        expires 30d;\n\
        add_header Cache-Control "public, immutable";\n\
    }\n\
    location = /index.html {\n\
        expires -1;\n\
        add_header Cache-Control "no-cache";\n\
    }\n\
}\n' > /etc/nginx/conf.d/default.conf

EXPOSE 8000

# Start both processes: config server in background, nginx in foreground
CMD config-server & nginx -g 'daemon off;'
