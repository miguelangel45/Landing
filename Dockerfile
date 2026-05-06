FROM nginx:alpine

# Copy the landing page and assets
COPY index.html /usr/share/nginx/html/index.html
COPY *.png *.svg /usr/share/nginx/html/

# Configure nginx to listen on port 8000 (Koyeb default)
RUN printf 'server {\n\
    listen 8000;\n\
    root /usr/share/nginx/html;\n\
    index index.html;\n\
    gzip on;\n\
    gzip_types text/html text/css application/javascript;\n\
}\n' > /etc/nginx/conf.d/default.conf

EXPOSE 8000

CMD ["nginx", "-g", "daemon off;"]
