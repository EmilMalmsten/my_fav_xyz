FROM node:16 AS builder
WORKDIR /app

# Define ARGs without default values
ARG VITE_API_URL_ARG
ARG VITE_IMG_URL_ARG

# Set the environment variables using ARG values
ENV VITE_API_URL=$VITE_API_URL_ARG
ENV VITE_IMG_URL=$VITE_IMG_URL_ARG

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Create a shared volume and copy the built files to it
FROM alpine:3.14 as fileserver
WORKDIR /app
COPY --from=builder /app/dist /app/dist

VOLUME ["/app/dist"]