import { useEffect, useState } from 'react'
import { Users, Loader2, Trash2, ShieldCheck, ShieldX, Crown, UserIcon } from 'lucide-react'
import { userService } from '@/services/endurance'
import type { User } from '@/types'
import { useAuth } from '@/contexts/AuthContext'
import toast from 'react-hot-toast'
import clsx from 'clsx'

export default function UsersPage() {
  const { userId } = useAuth()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)

  const fetchUsers = () => {
    setLoading(true)
    userService.getAll()
      .then((d) => setUsers(d.users))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchUsers() }, [])

  const handleToggleActive = async (u: User) => {
    try {
      await userService.update(u.id, { name: u.name, active: !u.active })
      toast.success(`Usuário ${u.active ? 'desativado' : 'ativado'}`)
      fetchUsers()
    } catch { /* handled */ }
  }

  const handleToggleRole = async (u: User) => {
    const newRole = u.role === 'admin' ? 'user' : 'admin'
    if (!confirm(`Alterar perfil de "${u.name}" para ${newRole === 'admin' ? 'Administrador' : 'Usuário'}?`)) return
    try {
      await userService.changeRole(u.id, newRole)
      toast.success(`Perfil alterado para ${newRole === 'admin' ? 'Administrador' : 'Usuário'}`)
      fetchUsers()
    } catch { /* handled */ }
  }

  const handleDelete = async (u: User) => {
    if (!confirm(`Deseja excluir "${u.name}"? Esta ação é irreversível.`)) return
    try {
      await userService.delete(u.id)
      toast.success('Usuário excluído')
      fetchUsers()
    } catch { /* handled */ }
  }

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-xl font-bold text-gray-900 dark:text-white">Usuários</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{users.length} usuário(s) cadastrado(s)</p>
      </div>

      {loading ? (
        <div className="flex justify-center py-16"><Loader2 className="w-6 h-6 animate-spin text-brand-500" /></div>
      ) : (
        <div className="card overflow-hidden p-0">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-100 dark:border-dark-border">
                {['Nome', 'E-mail', 'CPF', 'Perfil', 'Status', 'Ações'].map((h) => (
                  <th key={h} className="text-left text-xs font-semibold text-gray-500 dark:text-gray-400 px-4 py-3">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-50 dark:divide-dark-border">
              {users.map((u) => (
                <tr key={u.id} className="hover:bg-gray-50 dark:hover:bg-dark-muted/50 transition-colors">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2.5">
                      <div className="w-7 h-7 rounded-full bg-brand-100 dark:bg-brand-900/30 flex items-center justify-center text-xs font-bold text-brand-600 dark:text-brand-400">
                        {u.name.charAt(0).toUpperCase()}
                      </div>
                      <span className="font-medium text-gray-900 dark:text-white">{u.name}</span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-gray-500 dark:text-gray-400">{u.email}</td>
                  <td className="px-4 py-3 font-mono text-xs text-gray-500 dark:text-gray-400">{u.cpf}</td>
                  <td className="px-4 py-3">
                    <span className={clsx(u.role === 'admin' ? 'badge-critical' : 'badge-low')}>
                      {u.role === 'admin' ? 'Admin' : 'Usuário'}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span className={u.active ? 'badge-online' : 'badge-offline'}>
                      {u.active ? 'Ativo' : 'Inativo'}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-1">
                      {u.id !== userId && (
                        <>
                          <button
                            onClick={() => handleToggleRole(u)}
                            title={u.role === 'admin' ? 'Rebaixar para Usuário' : 'Promover a Admin'}
                            className="p-1.5 rounded hover:bg-gray-100 dark:hover:bg-dark-muted text-gray-400 hover:text-amber-500 dark:hover:text-amber-400 transition-colors"
                          >
                            {u.role === 'admin' ? <UserIcon className="w-3.5 h-3.5" /> : <Crown className="w-3.5 h-3.5" />}
                          </button>
                          <button
                            onClick={() => handleToggleActive(u)}
                            title={u.active ? 'Desativar' : 'Ativar'}
                            className="p-1.5 rounded hover:bg-gray-100 dark:hover:bg-dark-muted text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors"
                          >
                            {u.active ? <ShieldX className="w-3.5 h-3.5" /> : <ShieldCheck className="w-3.5 h-3.5" />}
                          </button>
                          <button
                            onClick={() => handleDelete(u)}
                            title="Excluir"
                            className="p-1.5 rounded hover:bg-red-50 dark:hover:bg-red-900/20 text-gray-400 hover:text-red-500 transition-colors"
                          >
                            <Trash2 className="w-3.5 h-3.5" />
                          </button>
                        </>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
