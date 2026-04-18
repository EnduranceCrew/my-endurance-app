import axios, { type AxiosError } from 'axios'
import toast from 'react-hot-toast'
import type { ApiEnvelope } from '@/types'

export const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
  timeout: 10_000,
})

// ── Interceptor de request: injeta o JWT ────────────────────────────────────

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('endurance_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// ── Interceptor de response: trata erros globalmente ────────────────────────

api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiEnvelope<unknown>>) => {
    const status = error.response?.status
    const message = error.response?.data?.error ?? 'Erro de conexão com o servidor'

    if (status === 401) {
      localStorage.removeItem('endurance_token')
      localStorage.removeItem('endurance_user')
      toast.error('Sessão expirada. Faça login novamente.')
      window.location.href = '/login'
    } else if (status === 403) {
      toast.error('Você não tem permissão para realizar esta ação.')
    } else if (status === 422 || status === 400) {
      toast.error(message)
    } else if (status && status >= 500) {
      toast.error('Erro interno do servidor. Tente novamente.')
    }

    return Promise.reject(error)
  },
)
