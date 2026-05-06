FROM nginx:alpine

# Copy the landing page and assets
COPY index.html /usr/share/nginx/html/index.html
COPY *.png *.svg *.webp /usr/share/nginx/html/

# Configure nginx to listen on port 8000 (Koyeb default)
RUN printf 'server {\n\
    listen 8000;\n\
    root /usr/share/nginx/html;\n\
    index index.html;\n\
    gzip on;\n\
    gzip_types text/html text/css application/javascript image/svg+xml;\n\
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

CMD ["nginx", "-g", "daemon off;"]
