# ---- Base Node ----
FROM node:20-alpine AS base
WORKDIR /app
COPY package.json pnpm-lock.yaml ./

# ---- Dependencies ----
FROM base AS dependencies
RUN npm install -g pnpm && pnpm install

# ---- Copy Files/Build ----
FROM dependencies AS build
COPY . .
RUN pnpm run build

# ---- Release ----
FROM nginx:alpine AS release
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
