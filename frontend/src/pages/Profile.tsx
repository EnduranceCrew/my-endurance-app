import { useState } from 'react'
import { User, Lock, Loader2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { userService } from '@/services/endurance'
import toast from 'react-hot-toast'

export default function Profile() {
  const { name, email, role, userId } = useAuth()
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({ current_password: '', new_password: '', confirm: '' })

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    if (form.new_password !== form.confirm) { toast.error('As senhas não coincidem'); return }
    if (form.new_password.length < 8) { toast.error('Senha deve ter ao menos 8 caracteres'); return }
    setLoading(true)
    try {
      await userService.changePassword({ current_password: form.current_password, new_password: form.new_password })
      toast.success('Senha alterada com sucesso!')
      setForm({ current_password: '', new_password: '', confirm: '' })
    } finally { setLoading(false) }
  }

  return (
    <div className="max-w-lg space-y-6">
      <h1 className="text-xl font-bold text-gray-900 dark:text-white">Meu perfil</h1>

      {/* Info */}
      <div className="card space-y-4">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-full bg-brand-100 dark:bg-brand-900/30 flex items-center justify-center">
            <User className="w-6 h-6 text-brand-600 dark:text-brand-400" />
          </div>
          <div>
            <p className="font-semibold text-gray-900 dark:text-white">{name}</p>
            <p className="text-sm text-gray-500 dark:text-gray-400">{email}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-500 dark:text-gray-400">Perfil:</span>
          <span className={`${role === 'admin' ? 'badge-critical' : 'badge-low'}`}>
            {role === 'admin' ? '⚡ Administrador' : '👤 Usuário'}
          </span>
        </div>
        <div>
          <span className="text-xs text-gray-500 dark:text-gray-400">ID:</span>
          <span className="ml-2 font-mono text-xs text-gray-400 dark:text-gray-600">{userId}</span>
        </div>
      </div>

      {/* Alterar senha */}
      <div className="card">
        <div className="flex items-center gap-2 mb-4">
          <Lock className="w-4 h-4 text-gray-500" />
          <h2 className="text-sm font-semibold text-gray-900 dark:text-white">Alterar senha</h2>
        </div>
        <form onSubmit={handleChangePassword} className="space-y-3">
          <div>
            <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Senha atual</label>
            <input type="password" className="input" value={form.current_password}
              onChange={(e) => setForm({ ...form, current_password: e.target.value })} required />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nova senha</label>
            <input type="password" className="input" minLength={8} value={form.new_password}
              onChange={(e) => setForm({ ...form, new_password: e.target.value })} required />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Confirmar nova senha</label>
            <input type="password" className="input" value={form.confirm}
              onChange={(e) => setForm({ ...form, confirm: e.target.value })} required />
          </div>
          <button type="submit" disabled={loading} className="btn-primary flex items-center gap-2">
            {loading && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
            Alterar senha
          </button>
        </form>
      </div>
    </div>
  )
}
