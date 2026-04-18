import { useEffect, useState } from 'react'
import { BellRing, Loader2, CheckCheck } from 'lucide-react'
import { alertService } from '@/services/endurance'
import type { Alert } from '@/types'
import AlertBadge from '@/components/AlertBadge'

export default function Alerts() {
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [loading, setLoading] = useState(true)
  const [onlyOpen, setOnlyOpen] = useState(true)

  const fetchAlerts = () => {
    setLoading(true)
    alertService.getAll(onlyOpen)
      .then((d) => setAlerts(d.alerts))
      .finally(() => setLoading(false))
  }

  useEffect(() => { fetchAlerts() }, [onlyOpen])

  return (
    <div className="space-y-5">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">Alertas</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{alerts.length} alerta(s)</p>
        </div>
        <button
          onClick={() => setOnlyOpen(!onlyOpen)}
          className={`flex items-center gap-2 text-xs px-3 py-2 rounded-lg border transition-colors ${
            onlyOpen
              ? 'border-brand-300 dark:border-brand-700 bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-400'
              : 'border-gray-200 dark:border-dark-border text-gray-500 hover:border-gray-300'
          }`}
        >
          <CheckCheck className="w-3.5 h-3.5" />
          {onlyOpen ? 'Mostrando abertos' : 'Mostrando todos'}
        </button>
      </div>

      {loading ? (
        <div className="flex justify-center py-16"><Loader2 className="w-6 h-6 animate-spin text-brand-500" /></div>
      ) : alerts.length === 0 ? (
        <div className="text-center py-16 text-gray-400 dark:text-gray-600">
          <BellRing className="w-10 h-10 mx-auto mb-3 opacity-30" />
          <p className="text-sm">{onlyOpen ? 'Nenhum alerta aberto 🎉' : 'Nenhum alerta registrado'}</p>
        </div>
      ) : (
        <div className="space-y-2">
          {alerts.map((a) => (
            <AlertBadge key={a.id} alert={a} onResolved={fetchAlerts} />
          ))}
        </div>
      )}
    </div>
  )
}
