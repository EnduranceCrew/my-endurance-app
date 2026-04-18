import { useState, type FormEvent } from 'react'
import { Navigate } from 'react-router-dom'
import { Eye, EyeOff, Cpu, Loader2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { useTheme } from '@/contexts/ThemeContext'
import ThemeToggle from '@/components/ThemeToggle'
import toast from 'react-hot-toast'

// ── Validação de CPF no frontend ─────────────────────────────────────────────

function formatCPF(value: string): string {
  const digits = value.replace(/\D/g, '').slice(0, 11)
  return digits
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d)/, '$1.$2')
    .replace(/(\d{3})(\d{1,2})$/, '$1-$2')
}

function isValidCPF(cpf: string): boolean {
  const d = cpf.replace(/\D/g, '')
  if (d.length !== 11 || /^(.)\1+$/.test(d)) return false
  const calc = (len: number) => {
    let sum = 0
    for (let i = 0; i < len; i++) sum += parseInt(d[i]) * (len + 1 - i)
    const rem = (sum * 10) % 11
    return rem >= 10 ? 0 : rem
  }
  return calc(9) === parseInt(d[9]) && calc(10) === parseInt(d[10])
}

type Tab = 'login' | 'register'

export default function Login() {
  const { isAuthenticated, login, register } = useAuth()
  const [tab, setTab] = useState<Tab>('login')
  const [loading, setLoading] = useState(false)
  const [showPass, setShowPass] = useState(false)

  // Campos
  const [email, setEmail]       = useState('')
  const [password, setPassword] = useState('')
  const [name, setName]         = useState('')
  const [cpf, setCPF]           = useState('')

  if (isAuthenticated) return <Navigate to="/dashboard" replace />

  const handleLogin = async (e: FormEvent) => {
    e.preventDefault()
    if (!email || !password) { toast.error('Preencha todos os campos'); return }
    setLoading(true)
    try { await login({ email, password }) }
    finally { setLoading(false) }
  }

  const handleRegister = async (e: FormEvent) => {
    e.preventDefault()
    if (!name || !email || !cpf || !password) { toast.error('Preencha todos os campos'); return }
    if (!isValidCPF(cpf)) { toast.error('CPF inválido'); return }
    if (password.length < 8) { toast.error('Senha deve ter ao menos 8 caracteres'); return }
    setLoading(true)
    try { await register({ name, email, cpf, password }) }
    finally { setLoading(false) }
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-bg flex flex-col">
      {/* Barra superior */}
      <div className="flex items-center justify-between px-6 py-4">
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 rounded-lg bg-brand-500 flex items-center justify-center">
            <Cpu className="w-3.5 h-3.5 text-white" />
          </div>
          <span className="font-bold text-sm text-gray-900 dark:text-white tracking-wide">ENDURANCE</span>
        </div>
        <ThemeToggle />
      </div>

      {/* Card central */}
      <div className="flex-1 flex items-center justify-center px-4 pb-8">
        <div className="w-full max-w-sm">
          {/* Header */}
          <div className="text-center mb-8">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
              {tab === 'login' ? 'Bem-vindo de volta' : 'Criar conta'}
            </h2>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
              {tab === 'login'
                ? 'Acesse o painel de monitoramento'
                : 'Preencha seus dados para começar'}
            </p>
          </div>

          {/* Tabs */}
          <div className="flex bg-gray-100 dark:bg-dark-surface rounded-lg p-1 mb-6">
            {(['login', 'register'] as Tab[]).map((t) => (
              <button
                key={t}
                onClick={() => setTab(t)}
                className={`flex-1 py-2 text-sm font-medium rounded-md transition-all duration-150 ${
                  tab === t
                    ? 'bg-white dark:bg-dark-card text-gray-900 dark:text-white shadow-sm'
                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'
                }`}
              >
                {t === 'login' ? 'Entrar' : 'Cadastrar'}
              </button>
            ))}
          </div>

          {/* Form Login */}
          {tab === 'login' && (
            <form onSubmit={handleLogin} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">E-mail</label>
                <input
                  type="email"
                  className="input"
                  placeholder="admin@endurance.dev"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  autoComplete="email"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Senha</label>
                <div className="relative">
                  <input
                    type={showPass ? 'text' : 'password'}
                    className="input pr-10"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    autoComplete="current-password"
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPass(!showPass)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
                  >
                    {showPass ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
              </div>
              <button type="submit" disabled={loading} className="btn-primary w-full mt-6 flex items-center justify-center gap-2">
                {loading && <Loader2 className="w-4 h-4 animate-spin" />}
                {loading ? 'Entrando...' : 'Entrar'}
              </button>
            </form>
          )}

          {/* Form Register */}
          {tab === 'register' && (
            <form onSubmit={handleRegister} className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nome completo</label>
                <input type="text" className="input" placeholder="João Silva" value={name} onChange={(e) => setName(e.target.value)} required />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">E-mail</label>
                <input type="email" className="input" placeholder="joao@exemplo.com" value={email} onChange={(e) => setEmail(e.target.value)} required />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">
                  CPF
                  {cpf && (
                    <span className={`ml-2 text-[10px] font-normal ${isValidCPF(cpf) ? 'text-emerald-500' : 'text-red-500'}`}>
                      {isValidCPF(cpf) ? '✓ válido' : '✗ inválido'}
                    </span>
                  )}
                </label>
                <input
                  type="text"
                  className="input font-mono"
                  placeholder="000.000.000-00"
                  value={cpf}
                  onChange={(e) => setCPF(formatCPF(e.target.value))}
                  maxLength={14}
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Senha</label>
                <div className="relative">
                  <input
                    type={showPass ? 'text' : 'password'}
                    className="input pr-10"
                    placeholder="Mínimo 8 caracteres"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    minLength={8}
                    required
                  />
                  <button type="button" onClick={() => setShowPass(!showPass)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                    {showPass ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                <div className="flex gap-1 mt-1.5">
                  {[6, 8, 10, 12].map((len) => (
                    <div key={len} className={`h-1 flex-1 rounded-full transition-colors ${password.length >= len ? 'bg-brand-500' : 'bg-gray-200 dark:bg-dark-muted'}`} />
                  ))}
                </div>
              </div>
              <button type="submit" disabled={loading} className="btn-primary w-full mt-6 flex items-center justify-center gap-2">
                {loading && <Loader2 className="w-4 h-4 animate-spin" />}
                {loading ? 'Criando conta...' : 'Criar conta'}
              </button>
            </form>
          )}

          <p className="text-center text-[11px] text-gray-400 dark:text-gray-600 mt-6">
            Endurance v1.0 · Monitoramento de Laboratórios de Informática
          </p>
        </div>
      </div>
    </div>
  )
}
