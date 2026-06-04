import './App.css'
import { useEffect, useState } from 'react'
import { Toaster } from 'react-hot-toast'
import CreateTaskPage from '@/pages/CreateTaskPage'
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

  const page = route === 'create-task' ? <CreateTaskPage /> : <HomePage />

  return (
    <main className="min-h-screen overflow-hidden bg-background text-foreground">
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 3200,
          style: {
            background: '#111614',
            border: '1px solid rgba(218, 232, 221, 0.16)',
            borderRadius: '8px',
            color: '#eef4ef',
            fontFamily: 'Geist Mono, ui-monospace, monospace',
            fontSize: '12px',
          },
          success: {
            iconTheme: {
              primary: '#72d69c',
              secondary: '#090c0b',
            },
          },
          error: {
            iconTheme: {
              primary: '#ff9f9f',
              secondary: '#090c0b',
            },
          },
        }}
      />
      {page}
    </main>
  )
}

export default App
