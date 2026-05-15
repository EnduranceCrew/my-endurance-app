# Endurance — Computer Lab Monitoring

> Full-stack system for monitoring public computer laboratories.  
> Dashboard with **dark** and **light** themes, JWT authentication, admin/user roles, CPF and email validation, pop-up notifications, and clean hexagonal architecture

---

## 🗂️ Project Structure

```
endurance/
├── backend/                     ← Go API (hexagonal architecture)
│   ├── cmd/api/main.go          ← Entry point + dependency injection
│   ├── config/                  ← Configuration and DB connection
│   ├── internal/
│   │   ├── domain/              ← ① DOMAIN: pure entities + repository interfaces
│   │   │   ├── user/
│   │   │   ├── lab/
│   │   │   ├── computer/
│   │   │   └── alert/
│   │   ├── application/         ← ② APPLICATION: use cases (business rules)
│   │   │   ├── auth/
│   │   │   ├── user/
│   │   │   ├── lab/
│   │   │   ├── computer/
│   │   │   ├── alert/
│   │   │   └── dashboard/
│   │   └── infrastructure/      ← ③ INFRASTRUCTURE: concrete adapters
│   │       ├── persistence/     ← GORM repositories (PostgreSQL)
│   │       ├── security/        ← JWT + bcrypt
│   │       └── http/            ← Gin handlers + middleware + router
│   └── pkg/                     ← Shared utilities (no business logic)
│       ├── apperrors/           ← Typed errors with HTTP codes
│       ├── response/            ← Standard JSON envelope
│       └── validator/           ← CPF validation (official algorithm) and email
└── frontend/                    ← SPA in TypeScript + React + Tailwind
    └── src/
        ├── contexts/            ← AuthContext (JWT) + ThemeContext (dark/light)
        ├── services/            ← Configured Axios + interceptors
        ├── components/          ← Layout, Sidebar, Navbar, StatsCard, LabCard…
        └── pages/               ← Login, Dashboard, Labs, LabDetail, Alerts, Users
```

---

## 🏗️ Hexagonal Architecture

```
┌─────────────────────────────────────────────┐
│              INFRASTRUCTURE                 │
│  ┌─────────────┐       ┌──────────────────┐ │
│  │  HTTP/Gin   │       │    PostgreSQL    │ │
│  │  (handlers) │       │   (GORM repos)   │ │
│  └──────┬──────┘       └────────┬─────────┘ │
│         │ primary port          │ secondary │
│  ┌──────▼──────────────────────▼─────────┐  │
│  │            APPLICATION                │  │
│  │   UseCase interfaces + implementations│  │
│  └──────────────────┬────────────────────┘  │
│                     │ domain ports          │
│  ┌──────────────────▼────────────────────┐  │
│  │               DOMAIN                  │  │
│  │  Pure entities · Interfaces (ports)   │  │
│  │  No framework dependency              │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

**Why hexagonal?**
- The domain is testable in isolation (no DB, no HTTP)
- Swapping PostgreSQL for another database = just reimplement the repositories
- Swapping Gin for another framework = just reimplement the handlers
- Use cases are the core: they receive and return DTOs, never framework objects

---

## ⚙️ Prerequisites

| Tool       | Minimum Version | Purpose        |
|------------|-----------------|----------------|
| Go         | 1.21            | Backend        |
| Node.js    | 18 LTS          | Frontend       |
| PostgreSQL | 14              | Database       |
| Git        | any             | Clone/version  |

### Quick Installation (Ubuntu/Debian)

```bash
# Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc

# Node.js via nvm (recommended)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
source ~/.bashrc
nvm install 20 && nvm use 20

# PostgreSQL
sudo apt update && sudo apt install -y postgresql postgresql-contrib
sudo systemctl start postgresql
```

---

## How to Run the Project from Scratch

### 1 · Database

```bash
# Enter PostgreSQL as superuser
sudo -u postgres psql

# Inside psql — create database and user:
CREATE DATABASE endurance;
CREATE USER endurance_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE endurance TO endurance_user;
\q
```

### 2 · Configure the backend environment

```bash
cd endurance/backend

# Copy the environment variables example file
cp .env.example .env
```

Open `.env` and edit the variables:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=endurance_user        # user created above
DB_PASSWORD=secure_password   # password created above
DB_NAME=endurance
JWT_SECRET=change_to_something_random_and_long_here
JWT_EXPIRATION_HOURS=24
CORS_ORIGINS=http://localhost:5173
```

> ⚠️ **Never commit `.env`** — it contains secrets. The `.env.example` is what goes into Git.

### 3 · Install backend dependencies

```bash
cd endurance/backend

# Downloads all Go modules declared in go.mod
go mod tidy

# Why run this?
# go.mod lists dependencies but doesn't download them automatically.
# `go mod tidy` downloads everything and generates go.sum (lock file).
```

### 4 · Run the backend

```bash
cd endurance/backend

# Compiles and starts the server
go run ./cmd/api/main.go

# You should see:
# [db] connected successfully!
# [migrate] tables synchronized
# 🔑 Default admin created: admin@endurance.dev / Admin@12345
# 🚀 Endurance running at http://localhost:8080
```

> The server creates tables automatically via GORM's AutoMigrate.  
> On first run, it creates the default admin user. **Change the password after your first login!**

### 5 · Install frontend dependencies

```bash
cd endurance/frontend

# Installs all dependencies listed in package.json
npm install

# Why npm install here?
# package.json declares the dependencies (React, Tailwind, Axios…)
# but node_modules/ doesn't exist yet. `npm install` creates that folder.
```

### 6 · Run the frontend

```bash
cd endurance/frontend

# Starts the Vite development server
npm run dev

# Vite opens at http://localhost:5173
# Proxy configured: /api → http://localhost:8080 (no manual CORS needed)
```

### 7 · Access the application

1. Open **http://localhost:5173**
2. Log in with the default admin credentials:
   - Email: `admin@endurance.dev`
   - Password: `Admin@12345`
3. **Change the password** under **My Profile → Change Password**

---

## 🔑 Credentials and Roles

| Role      | Permissions |
|-----------|-------------|
| **admin** | Everything: CRUD for labs, computers, users, alerts, dashboard |
| **user**  | View labs, computers, alerts; change own password |

---

## 📡 API Endpoints

### Auth (public)
| Method | Route | Description |
|--------|-------|-------------|
| POST | `/api/v1/auth/login` | Login → returns JWT |
| POST | `/api/v1/auth/register` | Registration → returns JWT |

### Dashboard (authenticated)
| Method | Route | Description |
|--------|-------|-------------|
| GET | `/api/v1/dashboard/stats` | General statistics |

### Laboratories
| Method | Route | Role |
|--------|-------|------|
| GET    | `/api/v1/labs` | all |
| GET    | `/api/v1/labs/:id` | all |
| POST   | `/api/v1/labs` | admin |
| PUT    | `/api/v1/labs/:id` | admin |
| DELETE | `/api/v1/labs/:id` | admin |
| GET    | `/api/v1/labs/:labId/computers` | all |
| GET    | `/api/v1/labs/:labId/alerts` | all |

### Computers
| Method | Route | Role |
|--------|-------|------|
| GET    | `/api/v1/computers` | all |
| POST   | `/api/v1/computers` | admin |
| PUT    | `/api/v1/computers/:id` | admin |
| PATCH  | `/api/v1/computers/:id/status` | all |
| DELETE | `/api/v1/computers/:id` | admin |

### Alerts
| Method | Route | Role |
|--------|-------|------|
| GET    | `/api/v1/alerts?open=true` | all |
| POST   | `/api/v1/alerts` | all |
| PATCH  | `/api/v1/alerts/:id/resolve` | admin |
| DELETE | `/api/v1/alerts/:id` | admin |

### Users
| Method | Route | Role |
|--------|-------|------|
| GET    | `/api/v1/users` | admin |
| PUT    | `/api/v1/users/:id` | admin |
| DELETE | `/api/v1/users/:id` | admin |
| POST   | `/api/v1/users/me/password` | authenticated |

---

## 🛠️ Useful Commands

### Backend

```bash
# Production build (generates binary)
go build -o endurance ./cmd/api/main.go

# Run the binary
./endurance

# Run tests
go test ./...

# Check for code issues
go vet ./...
```

### Frontend

```bash
# Production build
npm run build
# Generates the dist/ folder with static files

# Preview the production build
npm run preview

# Lint
npm run lint
```

---

## 🐳 Docker (optional)

```bash
# Backend
cat > backend/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o endurance ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/endurance .
EXPOSE 8080
CMD ["./endurance"]
EOF

# Frontend
cat > frontend/Dockerfile << 'EOF'
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json .
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
EOF

# docker-compose.yml at root
cat > docker-compose.yml << 'EOF'
version: '3.9'
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: endurance
      POSTGRES_USER: endurance_user
      POSTGRES_PASSWORD: secure_password
    ports: ["5432:5432"]
    volumes: [pgdata:/var/lib/postgresql/data]

  backend:
    build: ./backend
    ports: ["8080:8080"]
    env_file: ./backend/.env
    depends_on: [db]

  frontend:
    build: ./frontend
    ports: ["80:80"]
    depends_on: [backend]

volumes:
  pgdata:
EOF

# Start everything
docker-compose up --build
```

---

## 🎨 Themes

The system supports **dark** and **light** with a single click on the ☀️/🌙 icon in the navbar.  
The preference is saved in `localStorage` and respects the system's `prefers-color-scheme`.

- **Dark**: background `#0a0a0f`, cards `#16161f`, borders `#1e1e2a`
- **Light**: bright white, `gray-100` borders, subtle shadows
- **Accent**: blue `brand-500` (`#0ea5e9`) in both themes

---

## ✅ Implemented Validations

| Field    | Validation |
|----------|------------|
| CPF      | Official 2-digit verifier algorithm (frontend + backend) |
| Email    | Simplified RFC 5322 regex (frontend + backend) |
| Password | Minimum 8 characters + visual strength indicator on registration |
| Pop-ups  | `react-hot-toast` — success, error, warning, info with icons |

---

## 📦 Full Stack

| Layer        | Technology |
|--------------|-----------|
| Backend HTTP | Go + Gin |
| ORM          | GORM v2 |
| Database     | PostgreSQL |
| Auth         | JWT (golang-jwt/jwt v5) |
| Password hash | bcrypt (golang.org/x/crypto) |
| Frontend     | React 18 + TypeScript + Vite |
| Styling      | Tailwind CSS v3 |
| HTTP client  | Axios |
| Routing      | React Router v6 |
| Notifications | react-hot-toast |
| Icons        | Lucide React |
| Charts       | Recharts |

---

## 🔒 Security

- Passwords hashed with **bcrypt** (default cost = 12 rounds)
- JWT with configurable expiration (default 24h)
- Role-based access control (RBAC) middleware
- CORS restricted to origins configured in `.env`
- Soft-delete on all entities (data is never physically deleted)
- Error responses do not expose internal details in production (`GIN_MODE=release`)

---

## © Copyright

```
Copyright (c) 2025 Arthur Fialho. All rights reserved.
```

**Endurance — Computer Lab Monitoring System**  
Created and developed by **Arthur Fialho**.

This software and all associated source code, documentation, design assets, and related files are the exclusive intellectual property of Arthur Fialho, protected under applicable copyright laws and international treaties.

### Permissions & Restrictions

| Action | Status |
|--------|--------|
| Personal and educational use | ✅ Permitted |
| Modification for personal use | ✅ Permitted |
| Distribution of original or modified copies | ❌ Not permitted without written authorization |
| Commercial use | ❌ Not permitted without written authorization |
| Sublicensing | ❌ Not permitted |
| Removal or alteration of this copyright notice | ❌ Not permitted |

### Terms

Redistribution, reproduction, or use of this project — in whole or in part, in any form or by any means, electronic or mechanical — is strictly prohibited without the prior written permission of the copyright holder.

Any unauthorized use, copying, modification, merger, publication, distribution, sublicensing, or sale of copies of this software may constitute copyright infringement and may be subject to civil and/or criminal penalties under applicable law.

---

*© 2026 Arthur Fialho. All rights reserved. Unauthorized use is strictly prohibited.*
