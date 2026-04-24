import type { Lab } from '@/types'
import { MapPin, Users, Activity, Wrench, XCircle, Monitor, Pencil, Trash2 } from 'lucide-react'
import { Link } from 'react-router-dom'
import clsx from 'clsx'

const statusConfig = {
  active:      { label: 'Ativo',         icon: Activity, className: 'badge-online' },
  inactive:    { label: 'Inativo',        icon: XCircle,  className: 'badge-offline' },
  maintenance: { label: 'Manutenção',    icon: Wrench,   className: 'badge-idle' },
}

interface LabCardProps {
  lab: Lab
  onEdit?: (lab: Lab) => void
  onDelete?: (lab: Lab) => void
}

export default function LabCard({ lab, onEdit, onDelete }: LabCardProps) {
  const cfg = statusConfig[lab.status]
  const StatusIcon = cfg.icon
  const utilization = lab.computer_count > 0
    ? Math.round((lab.online_count / lab.computer_count) * 100)
    : 0

  return (
    <div className="relative group">
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
            <span><strong className="text-emerald-600 dark:text-emerald-400">{lab.online_count}</strong>/{lab.computer_count} online</span>
          </div>
        </div>

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

      {(onEdit || onDelete) && (
        <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
          {onEdit && (
            <button
              onClick={(e) => { e.preventDefault(); e.stopPropagation(); onEdit(lab) }}
              title="Editar"
              className="p-1 rounded bg-white dark:bg-dark-card shadow border border-gray-100 dark:border-dark-border text-gray-400 hover:text-brand-500 transition-colors"
            >
              <Pencil className="w-3 h-3" />
            </button>
          )}
          {onDelete && (
            <button
              onClick={(e) => { e.preventDefault(); e.stopPropagation(); onDelete(lab) }}
              title="Excluir"
              className="p-1 rounded bg-white dark:bg-dark-card shadow border border-gray-100 dark:border-dark-border text-gray-400 hover:text-red-500 transition-colors"
            >
              <Trash2 className="w-3 h-3" />
            </button>
          )}
        </div>
      )}
    </div>
  )
}
