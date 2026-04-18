import type { Alert } from '@/types'
import { AlertTriangle, Info, Wrench, Wifi, Zap, CheckCircle2 } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { alertService } from '@/services/endurance'
import toast from 'react-hot-toast'
import clsx from 'clsx'

const typeIcon = {
  offline:     Wifi,
  error:       AlertTriangle,
  maintenance: Wrench,
  overload:    Zap,
  info:        Info,
}

const severityClass = {
  low:      'badge-low',
  medium:   'badge-medium',
  high:     'badge-high',
  critical: 'badge-critical',
}

interface AlertBadgeProps {
  alert: Alert
  onResolved?: () => void
}

export default function AlertBadge({ alert, onResolved }: AlertBadgeProps) {
  const { isAdmin } = useAuth()
  const Icon = typeIcon[alert.type]

  const handleResolve = async () => {
    try {
      await alertService.resolve(alert.id)
      toast.success('Alerta resolvido!')
      onResolved?.()
    } catch { /* handled by interceptor */ }
  }

  return (
    <div className={clsx(
      'flex items-start gap-3 p-3 rounded-lg border transition-colors',
      alert.resolved
        ? 'bg-gray-50 dark:bg-dark-muted border-gray-100 dark:border-dark-border opacity-60'
        : 'bg-white dark:bg-dark-card border-gray-200 dark:border-dark-border'
    )}>
      <div className={clsx('mt-0.5 flex-shrink-0', severityClass[alert.severity])}>
        <Icon className="w-3.5 h-3.5" />
      </div>

      <div className="flex-1 min-w-0">
        <p className="text-sm text-gray-800 dark:text-gray-200 leading-snug">{alert.message}</p>
        <p className="text-[10px] text-gray-400 dark:text-gray-600 mt-1">
          {new Date(alert.created_at).toLocaleString('pt-BR')}
          {alert.resolved && alert.resolved_at && (
            <> · Resolvido em {new Date(alert.resolved_at).toLocaleString('pt-BR')}</>
          )}
        </p>
      </div>

      {!alert.resolved && isAdmin && (
        <button
          onClick={handleResolve}
          title="Marcar como resolvido"
          className="flex-shrink-0 p-1 rounded hover:bg-emerald-50 dark:hover:bg-emerald-900/20 text-gray-300 hover:text-emerald-500 transition-colors"
        >
          <CheckCircle2 className="w-4 h-4" />
        </button>
      )}
    </div>
  )
}
