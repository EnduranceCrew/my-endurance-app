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
import Computers from '@/pages/Computers'
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
