import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, Loader2, MapPin, Users, LayoutGrid, List } from 'lucide-react'
import { labService } from '@/services/endurance'
import type { Lab, Computer, Alert } from '@/types'
import ComputerGrid from '@/components/ComputerGrid'
import AlertBadge from '@/components/AlertBadge'

export default function LabDetail() {
  const { id } = useParams<{ id: string }>()
  const [lab, setLab] = useState<Lab | null>(null)
  const [computers, setComputers] = useState<Computer[]>([])
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [loading, setLoading] = useState(true)
  const [compact, setCompact] = useState(true)

  useEffect(() => {
    if (!id) return
    Promise.all([
      labService.getById(id),
      labService.getComputers(id),
      labService.getAlerts(id, false),
    ])
      .then(([l, c, a]) => { setLab(l); setComputers(c); setAlerts(a) })
      .finally(() => setLoading(false))
  }, [id])

  const refreshAlerts = () => labService.getAlerts(id!, false).then(setAlerts)

  if (loading) return <div className="flex justify-center py-16"><Loader2 className="w-6 h-6 animate-spin text-brand-500" /></div>
  if (!lab) return <div className="text-center py-16 text-gray-400">Laboratório não encontrado</div>

  const onlineCount = computers.filter((c) => c.status === 'online').length
  const openAlerts = alerts.filter((a) => !a.resolved)

  const statusColor = { active: 'bg-emerald-500', inactive: 'bg-gray-400', maintenance: 'bg-yellow-400' }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2">
        <Link to="/labs" className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors">
          <ArrowLeft className="w-4 h-4" />
        </Link>
        <span className="text-gray-300 dark:text-gray-600">/</span>
        <span className="text-sm text-gray-500 dark:text-gray-400">Laboratórios</span>
        <span className="text-gray-300 dark:text-gray-600">/</span>
        <span className="text-sm font-medium text-gray-900 dark:text-white">{lab.name}</span>
      </div>

      {/* Header do lab */}
      <div className="card">
        <div className="flex items-start justify-between gap-4">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <div className={`w-2 h-2 rounded-full ${statusColor[lab.status]}`} />
              <h1 className="text-lg font-bold text-gray-900 dark:text-white">{lab.name}</h1>
            </div>
            <div className="flex items-center gap-1.5 text-sm text-gray-500 dark:text-gray-400">
              <MapPin className="w-3.5 h-3.5" />
              {lab.location}
            </div>
            {lab.description && (
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">{lab.description}</p>
            )}
          </div>

          <div className="flex gap-6 text-center">
            <div>
              <p className="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{onlineCount}</p>
              <p className="text-[10px] text-gray-400 uppercase tracking-wide">Online</p>
            </div>
            <div>
              <p className="text-2xl font-bold text-gray-700 dark:text-gray-200">{computers.length}</p>
              <p className="text-[10px] text-gray-400 uppercase tracking-wide">Total</p>
            </div>
            <div>
              <p className={`text-2xl font-bold ${openAlerts.length ? 'text-red-500' : 'text-gray-700 dark:text-gray-200'}`}>
                {openAlerts.length}
              </p>
              <p className="text-[10px] text-gray-400 uppercase tracking-wide">Alertas</p>
            </div>
            <div className="flex items-center gap-1">
              <Users className="w-3.5 h-3.5 text-gray-400" />
              <div>
                <p className="text-2xl font-bold text-gray-700 dark:text-gray-200">{lab.capacity}</p>
                <p className="text-[10px] text-gray-400 uppercase tracking-wide">Capacidade</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Computadores */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-semibold text-gray-900 dark:text-white">
            Computadores ({computers.length})
          </h2>
          <div className="flex gap-1">
            <button onClick={() => setCompact(true)} className={`p-1.5 rounded ${compact ? 'bg-gray-100 dark:bg-dark-muted' : ''}`}>
              <LayoutGrid className="w-3.5 h-3.5 text-gray-500" />
            </button>
            <button onClick={() => setCompact(false)} className={`p-1.5 rounded ${!compact ? 'bg-gray-100 dark:bg-dark-muted' : ''}`}>
              <List className="w-3.5 h-3.5 text-gray-500" />
            </button>
          </div>
        </div>
        <ComputerGrid computers={computers} compact={compact} />
      </div>

      {/* Alertas */}
      {alerts.length > 0 && (
        <div>
          <h2 className="text-sm font-semibold text-gray-900 dark:text-white mb-3">
            Histórico de alertas ({alerts.length})
          </h2>
          <div className="space-y-2">
            {alerts.map((a) => (
              <AlertBadge key={a.id} alert={a} onResolved={refreshAlerts} />
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
