import type { Computer } from '@/types'
import clsx from 'clsx'
import { Monitor, Wifi, WifiOff, AlertTriangle, Clock, Trash2 } from 'lucide-react'

const statusConfig = {
  online:  { label: 'Online',       icon: Wifi,          className: 'badge-online',  dot: 'bg-emerald-500' },
  offline: { label: 'Offline',      icon: WifiOff,       className: 'badge-offline', dot: 'bg-gray-400' },
  error:   { label: 'Erro',         icon: AlertTriangle, className: 'badge-error',   dot: 'bg-red-500' },
  idle:    { label: 'Ocioso',       icon: Clock,         className: 'badge-idle',    dot: 'bg-yellow-400' },
}

interface ComputerGridProps {
  computers: Computer[]
  compact?: boolean
  onDelete?: (id: string) => void
}

export default function ComputerGrid({ computers, compact = false, onDelete }: ComputerGridProps) {
  if (computers.length === 0) {
    return (
      <div className="text-center py-12 text-gray-400 dark:text-gray-600">
        <Monitor className="w-10 h-10 mx-auto mb-3 opacity-40" />
        <p className="text-sm">Nenhum computador cadastrado</p>
      </div>
    )
  }

  if (compact) {
    return (
      <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-2">
        {computers.map((c) => {
          const cfg = statusConfig[c.status]
          return (
            <div
              key={c.id}
              title={`${c.hostname} (${c.ip_address}) — ${cfg.label}`}
              className="aspect-square rounded-lg bg-gray-50 dark:bg-dark-muted border border-gray-100 dark:border-dark-border flex flex-col items-center justify-center gap-1 cursor-default hover:scale-105 transition-transform"
            >
              <div className={clsx('w-1.5 h-1.5 rounded-full', cfg.dot)} />
              <span className="text-[9px] text-gray-500 dark:text-gray-400 font-mono leading-none truncate w-full text-center px-0.5">
                {c.hostname.replace(/[^0-9]/g, '') || c.hostname.slice(-4)}
              </span>
            </div>
          )
        })}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
      {computers.map((c) => {
        const cfg = statusConfig[c.status]
        const StatusIcon = cfg.icon
        return (
          <div
            key={c.id}
            className="card flex items-start gap-3 hover:border-brand-200 dark:hover:border-brand-800 transition-colors"
          >
            <div className="relative mt-0.5">
              <Monitor className="w-5 h-5 text-gray-400 dark:text-gray-600" />
              <span className={clsx('absolute -bottom-0.5 -right-0.5 w-2 h-2 rounded-full border border-white dark:border-dark-card', cfg.dot)} />
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center justify-between gap-2">
                <p className="text-sm font-medium text-gray-900 dark:text-white truncate font-mono">{c.hostname}</p>
                <span className={clsx(cfg.className, 'flex-shrink-0 text-[10px]')}>
                  <StatusIcon className="w-2.5 h-2.5" />
                  {cfg.label}
                </span>
              </div>
              <p className="text-xs text-gray-400 dark:text-gray-600 font-mono mt-0.5">{c.ip_address || '—'}</p>
              {(c.cpu || c.ram) && (
                <p className="text-[10px] text-gray-400 dark:text-gray-600 mt-1 truncate">
                  {[c.cpu, c.ram, c.os].filter(Boolean).join(' · ')}
                </p>
              )}
            </div>
            {onDelete && (
              <button
                onClick={() => onDelete(c.id)}
                title="Excluir"
                className="flex-shrink-0 p-1 rounded hover:bg-red-50 dark:hover:bg-red-900/20 text-gray-300 hover:text-red-500 transition-colors"
              >
                <Trash2 className="w-3.5 h-3.5" />
              </button>
            )}
          </div>
        )
      })}
    </div>
  )
}
