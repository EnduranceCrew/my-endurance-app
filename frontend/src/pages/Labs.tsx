import { useEffect, useState } from 'react'
import { Plus, Search, Loader2 } from 'lucide-react'
import { labService } from '@/services/endurance'
import type { Lab } from '@/types'
import LabCard from '@/components/LabCard'
import { useAuth } from '@/contexts/AuthContext'
import toast from 'react-hot-toast'

export default function Labs() {
  const { isAdmin } = useAuth()
  const [labs, setLabs] = useState<Lab[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [showModal, setShowModal] = useState(false)

  // Formulário de criação
  const [form, setForm] = useState({ name: '', location: '', capacity: 30, description: '' })
  const [creating, setCreating] = useState(false)

  const fetchLabs = () => {
    setLoading(true)
    labService.getAll()
      .then((d) => setLabs(d.labs))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchLabs() }, [])

  const filtered = labs.filter((l) =>
    l.name.toLowerCase().includes(search.toLowerCase()) ||
    l.location.toLowerCase().includes(search.toLowerCase())
  )

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name || !form.location) { toast.error('Preencha nome e localização'); return }
    setCreating(true)
    try {
      await labService.create(form)
      toast.success('Laboratório criado!')
      setShowModal(false)
      setForm({ name: '', location: '', capacity: 30, description: '' })
      fetchLabs()
    } finally { setCreating(false) }
  }

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">Laboratórios</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{labs.length} laboratório(s) cadastrado(s)</p>
        </div>
        {isAdmin && (
          <button onClick={() => setShowModal(true)} className="btn-primary flex items-center gap-2">
            <Plus className="w-4 h-4" /> Novo Lab
          </button>
        )}
      </div>

      {/* Busca */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          type="text"
          placeholder="Buscar laboratório..."
          className="input pl-9"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {loading ? (
        <div className="flex justify-center py-16"><Loader2 className="w-6 h-6 animate-spin text-brand-500" /></div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-16 text-gray-400 dark:text-gray-600">
          <p className="text-sm">Nenhum laboratório encontrado</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((lab) => <LabCard key={lab.id} lab={lab} />)}
        </div>
      )}

      {/* Modal criar lab */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
          <div className="card w-full max-w-md animate-slide-in">
            <h2 className="text-base font-semibold text-gray-900 dark:text-white mb-4">Novo laboratório</h2>
            <form onSubmit={handleCreate} className="space-y-3">
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Nome *</label>
                <input className="input" placeholder="Lab. Informática A" value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })} required />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Localização *</label>
                <input className="input" placeholder="Bloco B — Sala 201" value={form.location}
                  onChange={(e) => setForm({ ...form, location: e.target.value })} required />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Capacidade</label>
                <input type="number" className="input" min={1} max={200} value={form.capacity}
                  onChange={(e) => setForm({ ...form, capacity: Number(e.target.value) })} />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-700 dark:text-gray-300 mb-1.5">Descrição</label>
                <textarea className="input resize-none" rows={2} placeholder="Descrição opcional..."
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })} />
              </div>
              <div className="flex gap-2 pt-2">
                <button type="button" onClick={() => setShowModal(false)} className="btn-ghost flex-1">Cancelar</button>
                <button type="submit" disabled={creating} className="btn-primary flex-1 flex items-center justify-center gap-2">
                  {creating && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
                  Criar
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
