FROM node:20-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN if [ -f package-lock.json ]; then \
      npm ci --no-audit --fund=false; \
    else \
      npm install --no-audit --fund=false; \
    fi

FROM node:20-alpine AS dev
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NODE_ENV=production
RUN npm run build

FROM nginx:1.27-alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html

RUN rm -f /etc/nginx/conf.d/default.conf && \
    printf '%s\n' \
    'server {' \
    '  listen 80;' \
    '  server_name _;' \
    '  root /usr/share/nginx/html;' \
    '  index index.html;' \
    '  # Serve static files and fall back to index.html for client-side routing' \
    '  location / {' \
    '    try_files $uri $uri/ /index.html;' \
    '  }' \
    '  # Cache static assets aggressively (adjust as needed)' \
    '  location ~* \.(?:js|css|png|jpg|jpeg|gif|svg|ico|ttf|woff|woff2)$ {' \
    '    expires 1y;' \
    '    add_header Cache-Control "public, immutable";' \
    '    try_files $uri =404;' \
    '  }' \
    '}' \
    > /etc/nginx/conf.d/default.conf

EXPOSE 80
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD wget -qO- http://localhost/ >/dev/null 2>&1 || exit 1

CMD ["nginx", "-g", "daemon off;"]