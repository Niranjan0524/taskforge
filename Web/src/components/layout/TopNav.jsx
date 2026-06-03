import { Sparkles } from 'lucide-react'

function TopNav() {
  return (
    <nav className="topbar" aria-label="Main navigation">
      <a className="brand" href="#">
        <span className="brand-mark">
          <Sparkles size={18} strokeWidth={2.4} />
        </span>
        <span>TaskForge</span>
      </a>

      <div className="nav-links">
        <a href="#dashboard">Dashboard</a>
        <a href="#create-task">Create</a>
        <a href="#status">Status</a>
      </div>
    </nav>
  )
}

export default TopNav
