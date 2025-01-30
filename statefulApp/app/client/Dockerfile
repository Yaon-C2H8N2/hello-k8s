FROM node:22.0-alpine AS builder

WORKDIR /app

COPY ./package.json ./package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine as run

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80 443