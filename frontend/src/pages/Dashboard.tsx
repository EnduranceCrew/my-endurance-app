import { useEffect, useState } from 'react'
import { FlaskConical, Monitor, BellRing, Users, Zap, ShieldAlert, Wifi } from 'lucide-react'
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend,
} from 'recharts'
import StatsCard from '@/components/StatsCard'
import { dashboardService } from '@/services/endurance'
import type { DashboardStats } from '@/types'
import { useAuth } from '@/contexts/AuthContext'

// Dados simulados para o gráfico de linha (monitoramento ao longo do tempo)
const mockTimeline = [
  { time: '00:00', online: 42, offline: 18 },
  { time: '04:00', online: 12, offline: 48 },
  { time: '08:00', online: 55, offline: 5  },
  { time: '10:00', online: 60, offline: 0  },
  { time: '12:00', online: 58, offline: 2  },
  { time: '14:00', online: 62, offline: 0  },
  { time: '16:00', online: 57, offline: 5  },
  { time: '18:00', online: 45, offline: 15 },
  { time: '20:00', online: 30, offline: 30 },
  { time: '22:00', online: 20, offline: 40 },
]

const PIE_COLORS = ['#10b981', '#94a3b8', '#ef4444', '#f59e0b']

export default function Dashboard() {
  const { name, isAdmin } = useAuth()
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    dashboardService.getStats()
      .then(setStats)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return (
    <div className="flex items-center justify-center h-64">
      <div className="w-6 h-6 border-2 border-brand-500 border-t-transparent rounded-full animate-spin" />
    </div>
  )

  const pieData = stats ? [
    { name: 'Online',  value: stats.online_computers },
    { name: 'Offline', value: stats.offline_computers },
    { name: 'Erro',    value: stats.error_computers },
    { name: 'Ocioso',  value: stats.total_computers - stats.online_computers - stats.offline_computers - stats.error_computers },
  ].filter((d) => d.value > 0) : []

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-xl font-bold text-gray-900 dark:text-white">
          Olá, {name?.split(' ')[0]} 👋
        </h1>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
          {new Date().toLocaleDateString('pt-BR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
        </p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-5 gap-4">
        <StatsCard title="Laboratórios" value={stats?.total_labs ?? 0}
          icon={<FlaskConical className="w-5 h-5" />} color="blue"
          subtitle={`${stats?.active_labs} ativos`} />
        <StatsCard title="Computadores" value={stats?.total_computers ?? 0}
          icon={<Monitor className="w-5 h-5" />} color="purple"
          subtitle={`${stats?.online_computers} online`} />
        <StatsCard title="Online agora" value={stats?.online_computers ?? 0}
          icon={<Wifi className="w-5 h-5" />} color="green" />
        <StatsCard title="Alertas abertos" value={stats?.open_alerts ?? 0}
          icon={<BellRing className="w-5 h-5" />}
          color={stats?.open_alerts ? 'yellow' : 'green'} />
        <StatsCard title="Críticos" value={stats?.critical_alerts ?? 0}
          icon={<ShieldAlert className="w-5 h-5" />}
          color={stats?.critical_alerts ? 'red' : 'green'} />
      </div>

      {isAdmin && (
        <StatsCard title="Usuários cadastrados" value={stats?.total_users ?? 0}
          icon={<Users className="w-5 h-5" />} color="blue" />
      )}

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Timeline */}
        <div className="card lg:col-span-2">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-semibold text-gray-900 dark:text-white">Computadores online nas últimas 24h</h3>
            <span className="text-[10px] text-gray-400 dark:text-gray-600">Simulado</span>
          </div>
          <ResponsiveContainer width="100%" height={200}>
            <AreaChart data={mockTimeline}>
              <defs>
                <linearGradient id="gradOnline" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#0ea5e9" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#0ea5e9" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#1e1e2a" />
              <XAxis dataKey="time" tick={{ fontSize: 10, fill: '#6b7280' }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 10, fill: '#6b7280' }} axisLine={false} tickLine={false} />
              <Tooltip
                contentStyle={{ background: '#16161f', border: '1px solid #1e1e2a', borderRadius: 8, fontSize: 12 }}
                labelStyle={{ color: '#e5e7eb' }}
              />
              <Area type="monotone" dataKey="online" stroke="#0ea5e9" strokeWidth={2} fill="url(#gradOnline)" name="Online" />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        {/* Pie */}
        <div className="card">
          <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-4">Status geral</h3>
          {pieData.length > 0 ? (
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie data={pieData} cx="50%" cy="50%" innerRadius={50} outerRadius={75} paddingAngle={3} dataKey="value">
                  {pieData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip contentStyle={{ background: '#16161f', border: '1px solid #1e1e2a', borderRadius: 8, fontSize: 12 }} />
                <Legend iconType="circle" iconSize={8} wrapperStyle={{ fontSize: 11 }} />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <div className="flex items-center justify-center h-48 text-gray-400 text-sm">
              Sem dados ainda
            </div>
          )}
        </div>
      </div>

      {/* Alerta manutenção */}
      {stats?.maintenance_labs ? (
        <div className="flex items-center gap-3 p-4 rounded-xl bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800/40">
          <Zap className="w-4 h-4 text-yellow-600 dark:text-yellow-400 flex-shrink-0" />
          <p className="text-sm text-yellow-700 dark:text-yellow-300">
            <strong>{stats.maintenance_labs}</strong> laboratório(s) em manutenção no momento.
          </p>
        </div>
      ) : null}
    </div>
  )
}
