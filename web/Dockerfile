# Frontend build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY . .
RUN npm config set registry https://registry.npmmirror.com
RUN npm install && npm run build

# Nginx serve stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
