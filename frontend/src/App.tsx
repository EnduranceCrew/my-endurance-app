import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { AuthProvider } from '@/contexts/AuthContext'
import { ThemeProvider } from '@/contexts/ThemeContext'
import Layout from '@/components/Layout'
import Login from '@/pages/Login'
import Dashboard from '@/pages/Dashboard'
import Labs from '@/pages/Labs'
import LabDetail from '@/pages/LabDetail'
import Alerts from '@/pages/Alerts'
import UsersPage from '@/pages/Users'
import Profile from '@/pages/Profile'
import NotFound from '@/pages/NotFound'

export default function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <Toaster
            position="top-right"
            toastOptions={{
              duration: 4000,
              style: {
                background: 'var(--toast-bg, #fff)',
                color:      'var(--toast-color, #111)',
                fontSize:   '13px',
                fontWeight: '500',
                borderRadius: '10px',
                border: '1px solid rgba(0,0,0,0.06)',
                boxShadow: '0 4px 16px rgba(0,0,0,0.08)',
              },
              success: { iconTheme: { primary: '#10b981', secondary: '#fff' } },
              error:   { iconTheme: { primary: '#ef4444', secondary: '#fff' } },
            }}
          />
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route element={<Layout />}>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/labs"       element={<Labs />} />
              <Route path="/labs/:id"   element={<LabDetail />} />
              <Route path="/alerts"     element={<Alerts />} />
              <Route path="/computers"  element={<Computers />} />
              <Route path="/users"      element={<UsersPage />} />
              <Route path="/profile"    element={<Profile />} />
            </Route>
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  )
}

// Página de computadores inline (simples)
import { useEffect, useState } from 'react'
import { Monitor, Loader2 } from 'lucide-react'
import { computerService } from '@/services/endurance'
import type { Computer } from '@/types'
import ComputerGrid from '@/components/ComputerGrid'

function Computers() {
  const [computers, setComputers] = useState<Computer[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    computerService.getAll()
      .then((d) => setComputers(d.computers))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="space-y-5">
      <div className="flex items-center gap-2">
        <Monitor className="w-5 h-5 text-gray-500" />
        <h1 className="text-xl font-bold text-gray-900 dark:text-white">Computadores</h1>
        <span className="text-sm text-gray-400 dark:text-gray-600">({computers.length})</span>
      </div>
      {loading
        ? <div className="flex justify-center py-16"><Loader2 className="w-6 h-6 animate-spin text-brand-500" /></div>
        : <ComputerGrid computers={computers} />
      }
    </div>
  )
}
