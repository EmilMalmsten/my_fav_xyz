# Step 1: Build the application
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

# Step 2: Set up the production environment
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]