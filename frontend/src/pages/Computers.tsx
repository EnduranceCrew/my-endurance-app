import { useEffect, useState } from 'react'
import { Monitor, Loader2, Plus, Trash2 } from 'lucide-react'
import { computerService, labService } from '@/services/endurance'
import type { Computer, ComputerStatus, Lab } from '@/types'
import ComputerGrid from '@/components/ComputerGrid'
import { useAuth } from '@/contexts/AuthContext'
import toast from 'react-hot-toast'
import clsx from 'clsx'

const STATUS_TABS: { label: string; value: ComputerStatus | '' }[] = [
  { label: 'Todos', value: '' },
  { label: 'Online', value: 'online' },
  { label: 'Offline', value: 'offline' },
  { label: 'Erro', value: 'error' },
  { label: 'Ocioso', value: 'idle' },
]

const emptyForm = { lab_id: '', hostname: '', ip_address: '', mac_address: '', os: '', cpu: '', ram: '', storage: '' }

export default function Computers() {
  const { isAdmin } = useAuth()
  const [computers, setComputers] = useState<Computer[]>([])
  const [labs, setLabs] = useState<Lab[]>([])
  const [loading, setLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState<ComputerStatus | ''>('')
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)

  const fetchComputers = () => {
    setLoading(true)
    computerService.getAll(1, 200, statusFilter || undefined)
      .then((d) => setComputers(d.computers))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchComputers() }, [statusFilter])

  useEffect(() => {
    if (isAdmin) {
      labService.getAll().then((d) => setLabs(d.labs))
    }
  }, [isAdmin])

  const handleDelete = async (id: string) => {
    if (!confirm('Deseja excluir este computador?')) return
    try {
      await computerService.delete(id)
      toast.success('Computador excluído')
      fetchComputers()
    } catch { /* handled by interceptor */ }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.lab_id || !form.hostname) { toast.error('Preencha o laboratório e hostname'); return }
    setSaving(true)
    try {
      await computerService.create(form)
      toast.success('Computador adicionado!')
      setShowModal(false)
      setForm(emptyForm)
      fetchComputers()
    } finally { setSaving(false) }
  }

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Monitor className="w-5 h-5 text-gray-500" />
          <div>
            <h1 className="text-xl font-bold text-gray-900 dark:text-white">Computadores</h1>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{computers.length} computador(es)</p>
          </div>
        </div>
        {isAdmin && (
          <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
            <Plus className="w-4 h-4" /> Adicionar
          </button>
        )}
      </div>

      {/* Filtro por status */}
      <div className="flex gap-1.5 flex-wrap">
        {STATUS_TABS.map((tab) => (
          <button
            key={tab.value}
            onClick={() => setStatusFilter(tab.value)}
            className={clsx(
              'text-xs px-3 py-1.5 rounded-lg border transition-colors',
              statusFilter === tab.value
                ? 'border-brand-300 dark:border-brand-700 bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-400'
                : 'border-gray-200 dark:border-dark-border text-gray-500 hover:border-gray-300 dark:hover:border-gray-600'
            )}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="flex justify-center py-16"><Loader2 className="w-6 h-6 animate-spin text-brand-500" /></div>
      ) : (
        <ComputerGrid
          computers={computers}
          onDelete={isAdmin ? handleDelete : undefined}
        />
      )}

      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="card w-full max-w-md animate-slide-in">
            <h2 className="text-base font-semibold text-gray-900 dark:text-white mb-4">Novo computador</h2>
            <form onSubmit={handleCreate} className="space-y-3">
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Laboratório *</label>
                <select className="input" value={form.lab_id}
                  onChange={(e) => setForm({ ...form, lab_id: e.target.value })} required>
                  <option value="">Selecione...</option>
                  {labs.map((l) => (
                    <option key={l.id} value={l.id}>{l.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Hostname *</label>
                <input className="input" placeholder="pc-lab01" value={form.hostname}
                  onChange={(e) => setForm({ ...form, hostname: e.target.value })} required />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">IP</label>
                  <input className="input" placeholder="192.168.1.10" value={form.ip_address}
                    onChange={(e) => setForm({ ...form, ip_address: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">MAC</label>
                  <input className="input font-mono" placeholder="AA:BB:CC:DD:EE:FF" value={form.mac_address}
                    onChange={(e) => setForm({ ...form, mac_address: e.target.value })} />
                </div>
              </div>
              <div className="grid grid-cols-3 gap-3">
                <div>
                  <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">CPU</label>
                  <input className="input" placeholder="i5-12400" value={form.cpu}
                    onChange={(e) => setForm({ ...form, cpu: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">RAM</label>
                  <input className="input" placeholder="8GB" value={form.ram}
                    onChange={(e) => setForm({ ...form, ram: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Disco</label>
                  <input className="input" placeholder="256GB" value={form.storage}
                    onChange={(e) => setForm({ ...form, storage: e.target.value })} />
                </div>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">SO</label>
                <input className="input" placeholder="Windows 11" value={form.os}
                  onChange={(e) => setForm({ ...form, os: e.target.value })} />
              </div>
              <div className="flex gap-2 pt-2">
                <button type="button" onClick={() => setShowModal(false)} className="btn-ghost flex-1">Cancelar</button>
                <button type="submit" disabled={saving} className="btn-primary flex-1 flex items-center justify-center gap-2">
                  {saving && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
                  Adicionar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
