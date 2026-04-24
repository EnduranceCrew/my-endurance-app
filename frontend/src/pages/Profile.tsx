import { useState } from 'react'
import { User, Lock, Pencil, Loader2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { userService } from '@/services/endurance'
import toast from 'react-hot-toast'

export default function Profile() {
  const { name, email, role, userId } = useAuth()
  const [loadingName, setLoadingName] = useState(false)
  const [loadingPwd, setLoadingPwd] = useState(false)
  const [newName, setNewName] = useState(name ?? '')
  const [pwdForm, setPwdForm] = useState({ current_password: '', new_password: '', confirm: '' })

  const handleUpdateName = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newName.trim()) { toast.error('Nome não pode ser vazio'); return }
    if (newName.trim() === name) { toast('Nenhuma alteração detectada'); return }
    setLoadingName(true)
    try {
      await userService.update(userId!, { name: newName.trim() })
      toast.success('Nome atualizado!')
      // Força reload dos dados do usuário no contexto via /users/me
      const me = await userService.me()
      // Atualiza localStorage manualmente para refletir o novo nome sem re-login
      const raw = localStorage.getItem('endurance_user')
      if (raw) {
        const stored = JSON.parse(raw)
        stored.name = me.name
        localStorage.setItem('endurance_user', JSON.stringify(stored))
        window.location.reload()
      }
    } finally { setLoadingName(false) }
  }

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    if (pwdForm.new_password !== pwdForm.confirm) { toast.error('As senhas não coincidem'); return }
    if (pwdForm.new_password.length < 8) { toast.error('Senha deve ter ao menos 8 caracteres'); return }
    setLoadingPwd(true)
    try {
      await userService.changePassword({ current_password: pwdForm.current_password, new_password: pwdForm.new_password })
      toast.success('Senha alterada com sucesso!')
      setPwdForm({ current_password: '', new_password: '', confirm: '' })
    } finally { setLoadingPwd(false) }
  }

  return (
    <div className="max-w-lg space-y-6">
      <h1 className="text-xl font-bold text-gray-900 dark:text-white">Meu perfil</h1>

      {/* Info + editar nome */}
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

        {/* Editar nome */}
        <form onSubmit={handleUpdateName} className="flex gap-2 pt-1">
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nome</label>
            <input
              className="input"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              minLength={2}
              required
            />
          </div>
          <button type="submit" disabled={loadingName} className="btn-ghost self-end flex items-center gap-1.5">
            {loadingName ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Pencil className="w-3.5 h-3.5" />}
            Salvar
          </button>
        </form>
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
            <input type="password" className="input" value={pwdForm.current_password}
              onChange={(e) => setPwdForm({ ...pwdForm, current_password: e.target.value })} required />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nova senha</label>
            <input type="password" className="input" minLength={8} value={pwdForm.new_password}
              onChange={(e) => setPwdForm({ ...pwdForm, new_password: e.target.value })} required />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Confirmar nova senha</label>
            <input type="password" className="input" value={pwdForm.confirm}
              onChange={(e) => setPwdForm({ ...pwdForm, confirm: e.target.value })} required />
          </div>
          <button type="submit" disabled={loadingPwd} className="btn-primary flex items-center gap-2">
            {loadingPwd && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
            Alterar senha
          </button>
        </form>
      </div>
    </div>
  )
}
