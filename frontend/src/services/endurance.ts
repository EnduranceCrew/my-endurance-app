import { api } from './api'
import type { AuthToken, LoginInput, RegisterInput, DashboardStats, Lab, Computer, Alert, User } from '@/types'

// ── Auth ─────────────────────────────────────────────────────────────────────

export const authService = {
  login: (data: LoginInput) =>
    api.post<{ data: AuthToken }>('/auth/login', data).then((r) => r.data.data!),

  register: (data: RegisterInput) =>
    api.post<{ data: AuthToken }>('/auth/register', data).then((r) => r.data.data!),
}

// ── Dashboard ────────────────────────────────────────────────────────────────

export const dashboardService = {
  getStats: () =>
    api.get<{ data: DashboardStats }>('/dashboard/stats').then((r) => r.data.data!),
}

// ── Labs ─────────────────────────────────────────────────────────────────────

export const labService = {
  getAll: (page = 1, limit = 20) =>
    api.get<{ data: { labs: Lab[]; total: number } }>(`/labs?page=${page}&limit=${limit}`).then((r) => r.data.data!),

  getById: (id: string) =>
    api.get<{ data: Lab }>(`/labs/${id}`).then((r) => r.data.data!),

  create: (data: Partial<Lab>) =>
    api.post<{ data: Lab }>('/labs', data).then((r) => r.data.data!),

  update: (id: string, data: Partial<Lab>) =>
    api.put<{ data: Lab }>(`/labs/${id}`, data).then((r) => r.data.data!),

  delete: (id: string) => api.delete(`/labs/${id}`),

  getComputers: (labId: string) =>
    api.get<{ data: Computer[] }>(`/labs/${labId}/computers`).then((r) => r.data.data!),

  getAlerts: (labId: string, onlyOpen = true) =>
    api.get<{ data: Alert[] }>(`/labs/${labId}/alerts?open=${onlyOpen}`).then((r) => r.data.data!),
}

// ── Computers ────────────────────────────────────────────────────────────────

export const computerService = {
  getAll: (page = 1, limit = 20) =>
    api.get<{ data: { computers: Computer[]; total: number } }>(`/computers?page=${page}&limit=${limit}`).then((r) => r.data.data!),

  updateStatus: (id: string, status: Computer['status']) =>
    api.patch(`/computers/${id}/status`, { status }),

  create: (data: Partial<Computer>) =>
    api.post<{ data: Computer }>('/computers', data).then((r) => r.data.data!),

  delete: (id: string) => api.delete(`/computers/${id}`),
}

// ── Alerts ───────────────────────────────────────────────────────────────────

export const alertService = {
  getAll: (onlyOpen = true, page = 1, limit = 20) =>
    api.get<{ data: { alerts: Alert[]; total: number } }>(`/alerts?open=${onlyOpen}&page=${page}&limit=${limit}`).then((r) => r.data.data!),

  resolve: (id: string) => api.patch(`/alerts/${id}/resolve`),

  create: (data: Partial<Alert>) =>
    api.post<{ data: Alert }>('/alerts', data).then((r) => r.data.data!),

  delete: (id: string) => api.delete(`/alerts/${id}`),
}

// ── Users ────────────────────────────────────────────────────────────────────

export const userService = {
  getAll: (page = 1, limit = 20) =>
    api.get<{ data: { users: User[]; total: number } }>(`/users?page=${page}&limit=${limit}`).then((r) => r.data.data!),

  update: (id: string, data: Partial<User>) =>
    api.put<{ data: User }>(`/users/${id}`, data).then((r) => r.data.data!),

  delete: (id: string) => api.delete(`/users/${id}`),

  changePassword: (data: { current_password: string; new_password: string }) =>
    api.post('/users/me/password', data),
}
