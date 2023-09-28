
# 🚀 jl8n.dev

My personal developer webapp showcasing some of my development chops using some of my favorite technologies.

The current build can be viewed at [jl8n.dev](https://jl8n.dev).

## 📚 Tech Stack

### Backend
- **Go**
  - Go-chi API: Server-side events and HTTP requests.
- **MongoDB**
- **GraphQL**

### Frontend
- **Typescript**
- **Svelte & Sveltekit**



### Docker 🐳

Both the frontend and backend are dockerized as part of a CI/CD pipeline.


## 📦 Installation & Usage

### Run

Simply use `docker compose` to pull the latest Docker images and run the `web` and `server` containers:

```bash
docker compose up -d
```

Or using `docker run`:

```bash
docker run -d --name jl8n.dev-web -p 3001:80 ghcr.io/jl8n/jl8n.dev/web:latest
docker run -d --name jl8n.dev-server -p 3000:3000 ghcr.io/jl8n/jl8n.dev/server:latest
```

### Dev

Run the Vite development server:

```bash
cd web/
pnpm dev
```

Run the Go API service:

```bash
cd /server
go run *.go
```
