import './App.css'
import { useEffect, useState } from 'react'
import { Toaster } from 'react-hot-toast'
import CreateTaskPage from '@/pages/CreateTaskPage'
import DashboardPage from '@/pages/DashboardPage'
import HomePage from '@/pages/HomePage'

function getRouteFromHash() {
  return window.location.hash.replace('#', '') || 'home'
}

function App() {
  const [route, setRoute] = useState(getRouteFromHash)

  useEffect(() => {
    const handleHashChange = () => setRoute(getRouteFromHash())

    window.addEventListener('hashchange', handleHashChange)

    return () => window.removeEventListener('hashchange', handleHashChange)
  }, [])

  const pages = {
    dashboard: <DashboardPage />,
    'create-task': <CreateTaskPage />,
    home: <HomePage />,
  }

  const page = pages[route] || <HomePage />

  return (
    <main className="min-h-screen overflow-hidden bg-background text-foreground">
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 3200,
          style: {
            background: '#FCF8D8',
            border: '1px solid rgba(124, 125, 117, 0.28)',
            borderRadius: '8px',
            color: '#191A16',
            fontFamily: 'Geist Mono, ui-monospace, monospace',
            fontSize: '12px',
          },
          success: {
            iconTheme: {
              primary: '#DD700B',
              secondary: '#FCF8D8',
            },
          },
          error: {
            iconTheme: {
              primary: '#B3481B',
              secondary: '#FCF8D8',
            },
          },
        }}
      />
      {page}
    </main>
  )
}

export default App
