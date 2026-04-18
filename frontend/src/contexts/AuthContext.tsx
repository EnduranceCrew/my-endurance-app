import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import toast from 'react-hot-toast'
import type { AuthToken, LoginInput, RegisterInput, Role } from '@/types'
import { authService } from '@/services/endurance'

interface AuthState {
  token: string | null
  userId: string | null
  name: string | null
  email: string | null
  role: Role | null
}

interface AuthContextValue extends AuthState {
  isAuthenticated: boolean
  isAdmin: boolean
  login: (data: LoginInput) => Promise<void>
  register: (data: RegisterInput) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

const TOKEN_KEY = 'endurance_token'
const USER_KEY  = 'endurance_user'

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(() => {
    try {
      const raw = localStorage.getItem(USER_KEY)
      if (raw) return JSON.parse(raw) as AuthState
    } catch { /* ignore */ }
    return { token: null, userId: null, name: null, email: null, role: null }
  })

  // Sincroniza storage quando o estado muda
  useEffect(() => {
    if (state.token) {
      localStorage.setItem(TOKEN_KEY, state.token)
      localStorage.setItem(USER_KEY, JSON.stringify(state))
    }
  }, [state])

  const applyToken = (t: AuthToken) => {
    setState({
      token:  t.access_token,
      userId: t.user_id,
      name:   t.name,
      email:  t.email,
      role:   t.role,
    })
  }

  const login = async (data: LoginInput) => {
    const result = await authService.login(data)
    applyToken(result)
    toast.success(`Bem-vindo, ${result.name}! 👋`)
  }

  const register = async (data: RegisterInput) => {
    const result = await authService.register(data)
    applyToken(result)
    toast.success('Conta criada com sucesso!')
  }

  const logout = () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    setState({ token: null, userId: null, name: null, email: null, role: null })
    toast('Até logo!', { icon: '👋' })
  }

  return (
    <AuthContext.Provider value={{
      ...state,
      isAuthenticated: !!state.token,
      isAdmin: state.role === 'admin',
      login,
      register,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth deve ser usado dentro de <AuthProvider>')
  return ctx
}
