
Backend Go — Arquitetura Hexagonal

 3 camadas isoladas:
 ① DOMÍNIO     → entidades puras + interfaces (sem framework)
 ② APLICAÇÃO   → casos de uso com DTOs
 ③ INFRA       → GORM, Gin handlers, JWT, bcrypt

Frontend TypeScript

 - Login com tabs entrar/cadastrar, validação CPF em tempo real, indicador de força de senha
 - Dashboard com gráficos (Recharts), stats cards, alertas de manutenção
 - Tema dark/light com ThemeContext + Tailwind
 - JWT em AuthContext com interceptors Axios (401 → logout automático)
 - Pop-ups com react-hot-toast em todas as ações

Como rodar (resumo):

 # Terminal 1 — Backend
 cd backend && cp .env.example .env   # editar o .env
 go mod tidy
 go run ./cmd/api/main.go
 
 # Terminal 2 — Frontend  
 cd frontend && npm install
 npm run dev

Acesse http://localhost:5173 com:

 - 📧 admin@endurance.dev
 - 🔑 Admin@12345

O README.md tem todas as instruções detalhadas: cada comando explicado, endpoints, Docker, segurança, arquitetura.
