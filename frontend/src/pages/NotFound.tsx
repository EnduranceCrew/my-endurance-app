import { Link } from 'react-router-dom'
import { Home } from 'lucide-react'

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-dark-bg">
      <div className="text-center">
        <p className="text-8xl font-bold text-gray-200 dark:text-dark-border">404</p>
        <h1 className="text-xl font-semibold text-gray-900 dark:text-white mt-2">Página não encontrada</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1 mb-6">A rota que você acessou não existe.</p>
        <Link to="/dashboard" className="btn-primary inline-flex items-center gap-2">
          <Home className="w-4 h-4" /> Voltar ao início
        </Link>
      </div>
    </div>
  )
}
