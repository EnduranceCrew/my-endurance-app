import { Bell, LogOut, User } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import ThemeToggle from './ThemeToggle'

export default function Navbar() {
  const { name, role, logout } = useAuth()

  return (
    <header className="h-16 flex items-center justify-between px-6 bg-white dark:bg-dark-surface border-b border-gray-100 dark:border-dark-border flex-shrink-0">
      <div>
        <h1 className="text-sm font-semibold text-gray-900 dark:text-white">
          Sistema de Monitoramento
        </h1>
        <p className="text-xs text-gray-400 dark:text-gray-600">
          Laboratórios de Informática
        </p>
      </div>

      <div className="flex items-center gap-2">
        <ThemeToggle />

        <button className="relative p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-muted transition-colors">
          <Bell className="w-4 h-4 text-gray-500 dark:text-gray-400" />
        </button>

        <div className="flex items-center gap-2.5 pl-3 border-l border-gray-100 dark:border-dark-border ml-1">
          <div className="w-7 h-7 rounded-full bg-brand-100 dark:bg-brand-900/40 flex items-center justify-center">
            <User className="w-3.5 h-3.5 text-brand-600 dark:text-brand-400" />
          </div>
          <div className="hidden sm:block">
            <p className="text-xs font-medium text-gray-900 dark:text-white leading-none">{name}</p>
            <p className="text-[10px] text-gray-400 dark:text-gray-600 capitalize mt-0.5">{role}</p>
          </div>
          <button
            onClick={logout}
            className="p-1.5 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 text-gray-400 hover:text-red-500 dark:hover:text-red-400 transition-colors ml-1"
            title="Sair"
          >
            <LogOut className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </header>
  )
}
