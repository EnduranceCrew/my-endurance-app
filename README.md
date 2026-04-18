# Endurance — Monitoramento de Laboratórios de Informática

> Sistema full-stack para monitoramento de laboratórios públicos de informática.  
> Dashboard com temas **dark** e **light**, autenticação JWT, roles de admin/usuário, validação de CPF e e-mail, notificações pop-up e arquitetura hexagonal limpa.

---

## 🗂️ Estrutura do Projeto

```
endurance/
├── backend/                     ← API em Go (arquitetura hexagonal)
│   ├── cmd/api/main.go          ← Ponto de entrada + injeção de dependências
│   ├── config/                  ← Configurações e conexão com BD
│   ├── internal/
│   │   ├── domain/              ← ① DOMÍNIO: entidades puras + interfaces de repositório
│   │   │   ├── user/
│   │   │   ├── lab/
│   │   │   ├── computer/
│   │   │   └── alert/
│   │   ├── application/         ← ② APLICAÇÃO: casos de uso (regras de negócio)
│   │   │   ├── auth/
│   │   │   ├── user/
│   │   │   ├── lab/
│   │   │   ├── computer/
│   │   │   ├── alert/
│   │   │   └── dashboard/
│   │   └── infrastructure/      ← ③ INFRAESTRUTURA: adaptadores concretos
│   │       ├── persistence/     ← Repositórios GORM (PostgreSQL)
│   │       ├── security/        ← JWT + bcrypt
│   │       └── http/            ← Handlers Gin + middleware + roteador
│   └── pkg/                     ← Utilitários compartilhados (sem lógica de negócio)
│       ├── apperrors/           ← Erros tipados com código HTTP
│       ├── response/            ← Envelope JSON padrão
│       └── validator/           ← Validação de CPF (algoritmo oficial) e e-mail
└── frontend/                    ← SPA em TypeScript + React + Tailwind
    └── src/
        ├── contexts/            ← AuthContext (JWT) + ThemeContext (dark/light)
        ├── services/            ← Axios configurado + interceptors
        ├── components/          ← Layout, Sidebar, Navbar, StatsCard, LabCard…
        └── pages/               ← Login, Dashboard, Labs, LabDetail, Alerts, Users
```

---

## 🏗️ Arquitetura Hexagonal

```
┌─────────────────────────────────────────────┐
│              INFRAESTRUTURA                 │
│  ┌─────────────┐       ┌──────────────────┐ │
│  │  HTTP/Gin   │       │    PostgreSQL    │ │
│  │  (handlers) │       │   (GORM repos)   │ │
│  └──────┬──────┘       └────────┬─────────┘ │
│         │ primary port          │ secondary │
│  ┌──────▼──────────────────────▼─────────┐  │
│  │            APLICAÇÃO                  │  │
│  │   UseCase interfaces + implementações │  │
│  └──────────────────┬────────────────────┘  │
│                     │ domain ports          │
│  ┌──────────────────▼────────────────────┐  │
│  │              DOMÍNIO                  │  │
│  │  Entidades puras · Interfaces (ports) │  │
│  │  Sem dependência de framework         │  │
│  └───────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

**Por que hexagonal?**
- O domínio é testável isoladamente (sem BD, sem HTTP)
- Trocar PostgreSQL por outro banco = só reimplementar os repositórios
- Trocar Gin por outro framework = só reimplementar os handlers
- Use cases são o centro: recebem e retornam DTOs, nunca objetos de framework

---

## ⚙️ Pré-requisitos

| Ferramenta | Versão mínima | Para quê |
|------------|---------------|----------|
| Go         | 1.21          | Backend  |
| Node.js    | 18 LTS        | Frontend |
| PostgreSQL | 14            | Banco de dados |
| Git        | qualquer      | Clonar/versionar |

### Instalação rápida (Ubuntu/Debian)

```bash
# Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc

# Node.js via nvm (recomendado)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
source ~/.bashrc
nvm install 20 && nvm use 20

# PostgreSQL
sudo apt update && sudo apt install -y postgresql postgresql-contrib
sudo systemctl start postgresql
```

---

## 🚀 Como Rodar do Zero

### 1 · Banco de dados

```bash
# Entrar no PostgreSQL como superusuário
sudo -u postgres psql

# Dentro do psql — criar banco e usuário:
CREATE DATABASE endurance;
CREATE USER endurance_user WITH PASSWORD 'senha_segura';
GRANT ALL PRIVILEGES ON DATABASE endurance TO endurance_user;
\q
```

### 2 · Configurar o ambiente do backend

```bash
cd endurance/backend

# Copiar o arquivo de exemplo de variáveis de ambiente
cp .env.example .env
```

Abra `.env` e edite as variáveis:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=endurance_user        # usuário criado acima
DB_PASSWORD=senha_segura      # senha criada acima
DB_NAME=endurance
JWT_SECRET=mude_para_algo_aleatorio_e_longo_aqui
JWT_EXPIRATION_HOURS=24
CORS_ORIGINS=http://localhost:5173
```

> ⚠️ **Nunca commit o `.env`** — ele contém segredos. O `.env.example` é o que vai no Git.

### 3 · Instalar dependências do backend

```bash
cd endurance/backend

# Baixa todos os módulos Go declarados em go.mod
go mod tidy

# Por que rodar isso?
# O go.mod lista as dependências mas não as baixa automaticamente.
# `go mod tidy` baixa tudo e gera o go.sum (lock file).
```

### 4 · Rodar o backend

```bash
cd endurance/backend

# Compila e executa o servidor
go run ./cmd/api/main.go

# Você deve ver:
# [db] conectado com sucesso!
# [migrate] tabelas sincronizadas
# 🔑 Admin padrão criado: admin@endurance.dev / Admin@12345
# 🚀 Endurance rodando em http://localhost:8080
```

> O servidor cria as tabelas automaticamente via AutoMigrate do GORM.  
> Na primeira execução, cria o usuário admin padrão. **Troque a senha após o primeiro login!**

### 5 · Instalar dependências do frontend

```bash
cd endurance/frontend

# Instala todas as dependências listadas em package.json
npm install

# Por que npm install aqui?
# O package.json declara as dependências (React, Tailwind, Axios…)
# mas node_modules/ não existe ainda. `npm install` cria essa pasta.
```

### 6 · Rodar o frontend

```bash
cd endurance/frontend

# Inicia o servidor de desenvolvimento Vite
npm run dev

# Vite abre em http://localhost:5173
# Proxy configurado: /api → http://localhost:8080 (sem CORS manual)
```

### 7 · Acessar a aplicação

1. Abra **http://localhost:5173**
2. Login com as credenciais do admin padrão:
   - E-mail: `admin@endurance.dev`
   - Senha:  `Admin@12345`
3. **Troque a senha** em **Meu Perfil → Alterar senha**

---

## 🔑 Credenciais e Roles

| Role  | Permissões |
|-------|-----------|
| **admin** | Tudo: CRUD de labs, computadores, usuários, alertas, dashboard |
| **user**  | Visualizar labs, computadores, alertas; alterar própria senha |

---

## 📡 Endpoints da API

### Auth (público)
| Método | Rota | Descrição |
|--------|------|-----------|
| POST | `/api/v1/auth/login` | Login → retorna JWT |
| POST | `/api/v1/auth/register` | Cadastro → retorna JWT |

### Dashboard (autenticado)
| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/api/v1/dashboard/stats` | Estatísticas gerais |

### Laboratórios
| Método | Rota | Role |
|--------|------|------|
| GET    | `/api/v1/labs` | todos |
| GET    | `/api/v1/labs/:id` | todos |
| POST   | `/api/v1/labs` | admin |
| PUT    | `/api/v1/labs/:id` | admin |
| DELETE | `/api/v1/labs/:id` | admin |
| GET    | `/api/v1/labs/:labId/computers` | todos |
| GET    | `/api/v1/labs/:labId/alerts` | todos |

### Computadores
| Método | Rota | Role |
|--------|------|------|
| GET    | `/api/v1/computers` | todos |
| POST   | `/api/v1/computers` | admin |
| PUT    | `/api/v1/computers/:id` | admin |
| PATCH  | `/api/v1/computers/:id/status` | todos |
| DELETE | `/api/v1/computers/:id` | admin |

### Alertas
| Método | Rota | Role |
|--------|------|------|
| GET    | `/api/v1/alerts?open=true` | todos |
| POST   | `/api/v1/alerts` | todos |
| PATCH  | `/api/v1/alerts/:id/resolve` | admin |
| DELETE | `/api/v1/alerts/:id` | admin |

### Usuários
| Método | Rota | Role |
|--------|------|------|
| GET    | `/api/v1/users` | admin |
| PUT    | `/api/v1/users/:id` | admin |
| DELETE | `/api/v1/users/:id` | admin |
| POST   | `/api/v1/users/me/password` | autenticado |

---

## 🛠️ Comandos Úteis

### Backend

```bash
# Build para produção (gera binário)
go build -o endurance ./cmd/api/main.go

# Executar o binário
./endurance

# Rodar testes
go test ./...

# Verificar problemas de código
go vet ./...
```

### Frontend

```bash
# Build para produção
npm run build
# Gera a pasta dist/ com os arquivos estáticos

# Preview do build de produção
npm run preview

# Lint
npm run lint
```

---

## 🐳 Docker (opcional)

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

# docker-compose.yml na raiz
cat > docker-compose.yml << 'EOF'
version: '3.9'
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: endurance
      POSTGRES_USER: endurance_user
      POSTGRES_PASSWORD: senha_segura
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

# Subir tudo
docker-compose up --build
```

---

## 🎨 Temas

O sistema suporta **dark** e **light** com um clique no ícone ☀️/🌙 no navbar.  
A preferência é salva no `localStorage` e respeita `prefers-color-scheme` do sistema.

- **Dark**: fundo `#0a0a0f`, cards `#16161f`, bordas `#1e1e2a`
- **Light**: branco brilhante, bordas `gray-100`, sombras sutis
- **Detalhes**: azul `brand-500` (`#0ea5e9`) em ambos os temas

---

## ✅ Validações implementadas

| Campo | Validação |
|-------|-----------|
| CPF   | Algoritmo dos 2 dígitos verificadores (frontend + backend) |
| E-mail | Regex RFC 5322 simplificado (frontend + backend) |
| Senha | Mínimo 8 caracteres + indicador de força visual no cadastro |
| Pop-ups | `react-hot-toast` — sucesso, erro, aviso, info com ícones |

---

## 📦 Stack completa

| Camada | Tecnologia |
|--------|-----------|
| Backend HTTP | Go + Gin |
| ORM | GORM v2 |
| Banco de dados | PostgreSQL |
| Autenticação | JWT (golang-jwt/jwt v5) |
| Hash de senhas | bcrypt (golang.org/x/crypto) |
| Frontend | React 18 + TypeScript + Vite |
| Estilização | Tailwind CSS v3 |
| HTTP client | Axios |
| Roteamento | React Router v6 |
| Notificações | react-hot-toast |
| Ícones | Lucide React |
| Gráficos | Recharts |

---

## 🔒 Segurança

- Senhas hasheadas com **bcrypt** (custo padrão = 12 rounds)
- JWT com expiração configurável (padrão 24h)
- Middleware de role-based access control (RBAC)
- CORS restrito às origens configuradas em `.env`
- Soft-delete em todas as entidades (dados nunca apagados fisicamente)
- Respostas de erro não expõem detalhes internos em produção (`GIN_MODE=release`)

---
