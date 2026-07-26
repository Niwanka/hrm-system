# ==========================================
# STAGE 1: Build Vue Production Bundle
# ==========================================
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package management files
COPY package*.json ./
RUN npm install

# Copy frontend source code
COPY . .

# Build production bundle to /app/dist
RUN npm run build

# ==========================================
# STAGE 2: Nginx Web Server Container
# ==========================================
FROM nginx:alpine

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy production static assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
