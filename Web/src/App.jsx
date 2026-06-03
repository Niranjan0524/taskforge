import './App.css'
import { useEffect, useState } from 'react'
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
      {page}
    </main>
  )
}

export default App
