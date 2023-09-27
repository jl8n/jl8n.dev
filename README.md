
# 🚀 jl8n.dev

My personal developer webapp showcasing some of my development chops using some of my favorite technologies.

## 📚 Tech Stack

### Backend
- **Go** 🐹
  - Go-chi API: Server-side events and HTTP requests.
- **MongoDB** 🍃
- **GraphQL** 📈

### Frontend
- **Typescript** 🟦
- **Svelte & Sveltekit** 🟧


### Docker 🐳
- **Nginx**: Hosts the frontend inside the Docker image.

## 📦 Installation & Usage

### Frontend

#### Dev

Start the development server

```bash
cd web/
pnpm dev
```

#### Run

Pull the Docker image of the latest build

```bash
docker pull ghcr.io/jl8n/jl8n.dev/web:latest
```

### Backend


#### Dev

To run for development:

```bash
cd /server
go run *.go
```

#### Run

Pull the Docker image of the latest build:

```bash
docker pull ghcr.io/jl8n/jl8n.dev/web:latest
```
