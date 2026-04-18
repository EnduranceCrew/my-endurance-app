import type { ReactNode } from 'react'
import clsx from 'clsx'

interface StatsCardProps {
  title: string
  value: number | string
  icon: ReactNode
  trend?: { value: number; label: string }
  color?: 'blue' | 'green' | 'red' | 'yellow' | 'purple'
  subtitle?: string
}

const colorMap = {
  blue:   'bg-brand-50 dark:bg-brand-900/20 text-brand-600 dark:text-brand-400',
  green:  'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400',
  red:    'bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400',
  yellow: 'bg-yellow-50 dark:bg-yellow-900/20 text-yellow-600 dark:text-yellow-400',
  purple: 'bg-purple-50 dark:bg-purple-900/20 text-purple-600 dark:text-purple-400',
}

export default function StatsCard({ title, value, icon, color = 'blue', subtitle, trend }: StatsCardProps) {
  return (
    <div className="card flex items-start gap-4">
      <div className={clsx('w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0', colorMap[color])}>
        {icon}
      </div>

      <div className="flex-1 min-w-0">
        <p className="text-xs font-medium text-gray-500 dark:text-gray-400 truncate">{title}</p>
        <p className="text-2xl font-bold text-gray-900 dark:text-white mt-0.5 leading-none">{value}</p>
        {subtitle && (
          <p className="text-xs text-gray-400 dark:text-gray-600 mt-1">{subtitle}</p>
        )}
        {trend && (
          <p className={clsx('text-xs mt-1 font-medium', trend.value >= 0 ? 'text-emerald-500' : 'text-red-500')}>
            {trend.value >= 0 ? '▲' : '▼'} {Math.abs(trend.value)}% {trend.label}
          </p>
        )}
      </div>
    </div>
  )
}
