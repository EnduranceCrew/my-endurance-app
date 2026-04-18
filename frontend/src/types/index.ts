// ── Entidades de domínio ─────────────────────────────────────────────────────

export type Role = 'admin' | 'user'

export interface User {
  id: string
  name: string
  email: string
  cpf: string
  role: Role
  active: boolean
  created_at: string
}

export type LabStatus = 'active' | 'inactive' | 'maintenance'

export interface Lab {
  id: string
  name: string
  location: string
  capacity: number
  status: LabStatus
  description: string
  responsible_id?: string
  created_at: string
  updated_at: string
}

export type ComputerStatus = 'online' | 'offline' | 'error' | 'idle'

export interface Computer {
  id: string
  lab_id: string
  hostname: string
  ip_address: string
  mac_address: string
  status: ComputerStatus
  os: string
  cpu: string
  ram: string
  storage: string
  last_seen?: string
  created_at: string
}

export type AlertSeverity = 'low' | 'medium' | 'high' | 'critical'
export type AlertType = 'offline' | 'error' | 'maintenance' | 'overload' | 'info'

export interface Alert {
  id: string
  lab_id: string
  computer_id?: string
  type: AlertType
  severity: AlertSeverity
  message: string
  resolved: boolean
  resolved_at?: string
  created_at: string
}

// ── DTOs de API ──────────────────────────────────────────────────────────────

export interface ApiEnvelope<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

export interface PaginatedList<T> {
  total: number
  page: number
  limit: number
  items: T[]
}

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  name: string
  email: string
  cpf: string
  password: string
}

export interface AuthToken {
  access_token: string
  token_type: string
  user_id: string
  name: string
  email: string
  role: Role
}

export interface DashboardStats {
  total_labs: number
  active_labs: number
  maintenance_labs: number
  total_computers: number
  online_computers: number
  offline_computers: number
  error_computers: number
  total_users: number
  open_alerts: number
  critical_alerts: number
}
