import type { Lab } from '@/types'
import { MapPin, Users, Activity, Wrench, XCircle } from 'lucide-react'
import { Link } from 'react-router-dom'
import clsx from 'clsx'

const statusConfig = {
  active:      { label: 'Ativo',         icon: Activity, className: 'badge-online' },
  inactive:    { label: 'Inativo',        icon: XCircle,  className: 'badge-offline' },
  maintenance: { label: 'Manutenção',    icon: Wrench,   className: 'badge-idle' },
}

interface LabCardProps {
  lab: Lab
  computerCount?: number
  onlineCount?: number
}

export default function LabCard({ lab, computerCount = 0, onlineCount = 0 }: LabCardProps) {
  const cfg = statusConfig[lab.status]
  const StatusIcon = cfg.icon
  const utilization = computerCount > 0 ? Math.round((onlineCount / computerCount) * 100) : 0

  return (
    <Link
      to={`/labs/${lab.id}`}
      className="card hover:border-brand-200 dark:hover:border-brand-800 hover:shadow-md transition-all duration-200 block"
    >
      <div className="flex items-start justify-between mb-3">
        <h3 className="font-semibold text-gray-900 dark:text-white text-sm truncate pr-2">{lab.name}</h3>
        <span className={clsx(cfg.className, 'flex-shrink-0')}>
          <StatusIcon className="w-3 h-3" />
          {cfg.label}
        </span>
      </div>

      <div className="flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400 mb-4">
        <MapPin className="w-3.5 h-3.5 flex-shrink-0" />
        <span className="truncate">{lab.location}</span>
      </div>

      <div className="flex items-center justify-between text-xs mb-3">
        <div className="flex items-center gap-1 text-gray-500 dark:text-gray-400">
          <Users className="w-3.5 h-3.5" />
          <span>Capacidade: <strong className="text-gray-700 dark:text-gray-200">{lab.capacity}</strong></span>
        </div>
        <div className="flex items-center gap-1 text-gray-500 dark:text-gray-400">
          <Monitor className="w-3.5 h-3.5" />
          <span><strong className="text-emerald-600 dark:text-emerald-400">{onlineCount}</strong>/{computerCount} online</span>
        </div>
      </div>

      {/* Barra de utilização */}
      <div className="w-full h-1 bg-gray-100 dark:bg-dark-muted rounded-full overflow-hidden">
        <div
          className={clsx(
            'h-full rounded-full transition-all duration-500',
            utilization >= 80 ? 'bg-emerald-500' :
            utilization >= 40 ? 'bg-brand-500' : 'bg-gray-300 dark:bg-gray-600'
          )}
          style={{ width: `${utilization}%` }}
        />
      </div>
      <p className="text-[10px] text-gray-400 dark:text-gray-600 mt-1">{utilization}% em uso</p>
    </Link>
  )
}

// Import necessário no topo mas missing — adicionando inline
function Monitor({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={2}>
      <rect x="2" y="3" width="20" height="14" rx="2" />
      <path d="M8 21h8M12 17v4" />
    </svg>
  )
}
