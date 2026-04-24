import { AlertTriangle } from 'lucide-react'

interface PageStateProps {
  loading: boolean
  error?: Error | null
  children: React.ReactNode
}

/**
 * Componente utilitário que abstrai os estados de loading e erro de páginas.
 * Uso: <PageState loading={loading} error={error}>{conteúdo}</PageState>
 */
export default function PageState({ loading, error, children }: PageStateProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="w-6 h-6 border-2 border-brand-500 border-t-transparent rounded-full animate-spin" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-3 text-red-500 dark:text-red-400">
        <AlertTriangle className="w-8 h-8" />
        <p className="text-sm font-medium">Erro ao carregar dados</p>
        <p className="text-xs text-gray-500 dark:text-gray-400">{error.message}</p>
      </div>
    )
  }

  return <>{children}</>
}
