import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  FlaskConical,
  Monitor,
  BellRing,
  Users,
  Cpu,
  ShieldCheck,
} from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import clsx from 'clsx'

const navItems = [
  { to: '/dashboard',  label: 'Dashboard',     icon: LayoutDashboard, adminOnly: false },
  { to: '/labs',       label: 'Laboratórios',  icon: FlaskConical,    adminOnly: false },
  { to: '/computers',  label: 'Computadores',  icon: Monitor,         adminOnly: false },
  { to: '/alerts',     label: 'Alertas',       icon: BellRing,        adminOnly: false },
  { to: '/users',      label: 'Usuários',      icon: Users,           adminOnly: true  },
]

export default function Sidebar() {
  const { isAdmin } = useAuth()

  return (
    <aside className="hidden md:flex w-60 flex-shrink-0 flex-col bg-white dark:bg-dark-surface border-r border-gray-100 dark:border-dark-border">
      {/* Logo */}
      <div className="flex items-center gap-2.5 px-5 h-16 border-b border-gray-100 dark:border-dark-border">
        <div className="w-8 h-8 rounded-lg bg-brand-500 flex items-center justify-center">
          <Cpu className="w-4 h-4 text-white" />
        </div>
        <div>
          <span className="font-bold text-sm tracking-wide text-gray-900 dark:text-white">ENDURANCE</span>
          <p className="text-[10px] text-gray-400 dark:text-gray-600 uppercase tracking-widest">Monitor</p>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
        <p className="px-3 mb-2 text-[10px] font-semibold uppercase tracking-widest text-gray-400 dark:text-gray-600">
          Navegação
        </p>
        {navItems
          .filter((item) => !item.adminOnly || isAdmin)
          .map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                clsx('sidebar-link', isActive && 'active')
              }
            >
              <item.icon className="w-4 h-4 flex-shrink-0" />
              {item.label}
            </NavLink>
          ))}
      </nav>

      {/* Admin badge */}
      {isAdmin && (
        <div className="px-5 py-3 border-t border-gray-100 dark:border-dark-border">
          <div className="flex items-center gap-2 text-xs text-brand-600 dark:text-brand-400">
            <ShieldCheck className="w-3.5 h-3.5" />
            <span className="font-medium">Modo Administrador</span>
          </div>
        </div>
      )}
    </aside>
  )
}
