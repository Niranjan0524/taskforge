import './App.css'
import {
  Activity,
  ArrowRight,
  BarChart3,
  CheckCircle2,
  Clock3,
  FileJson,
  Gauge,
  Plus,
  ShieldCheck,
  Sparkles,
  Workflow,
} from 'lucide-react'

const actions = [
  {
    title: 'Open Dashboard',
    description: 'Watch queued, running, completed, and failed jobs from one control room.',
    icon: BarChart3,
    href: '#dashboard',
    primary: true,
  },
  {
    title: 'Create Task',
    description: 'Submit email jobs, report generation, or custom payloads into the queue.',
    icon: Plus,
    href: '#create-task',
  },
  {
    title: 'Check Status',
    description: 'Look up task state by ID and inspect retry progress without digging through logs.',
    icon: Activity,
    href: '#status',
  },
]

const stats = [
  { label: 'Queue states', value: '4', detail: 'pending to failed' },
  { label: 'Worker pool', value: '3', detail: 'parallel executors' },
  { label: 'Default retries', value: '3', detail: 'per task fallback' },
]

const workflow = [
  {
    title: 'Create',
    text: 'Send a task type, priority, retry limit, and JSON payload.',
    icon: FileJson,
  },
  {
    title: 'Process',
    text: 'Workers pick from Redis and move tasks into running state.',
    icon: Workflow,
  },
  {
    title: 'Resolve',
    text: 'Completed and failed jobs stay visible for review.',
    icon: CheckCircle2,
  },
]

function App() {
  return (
    <main className="min-h-screen overflow-hidden bg-background text-foreground">
      <section className="home-shell">
        <nav className="topbar" aria-label="Main navigation">
          <a className="brand" href="/">
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

        <div className="hero-grid">
          <div className="hero-copy">
            <div className="eyebrow">
              <Gauge size={16} />
              Redis backed task orchestration
            </div>

            <h1>Forge, track, and finish background work with clarity.</h1>

            <p className="hero-description">
              TaskForge gives your worker queue a clean command surface: create
              structured tasks, monitor live status, inspect payloads, and keep
              execution history easy to understand.
            </p>

            <div className="hero-actions">
              <a className="button button-primary" href="#dashboard">
                Go to Dashboard
                <ArrowRight size={18} />
              </a>
              <a className="button button-secondary" href="#create-task">
                Create Task
                <Plus size={18} />
              </a>
            </div>
          </div>

          <aside className="console-panel" aria-label="Task preview">
            <div className="console-header">
              <span>taskforge.queue</span>
              <div className="status-pill">
                <span></span>
                online
              </div>
            </div>

            <div className="task-card active">
              <div>
                <p>generate_report</p>
                <span>priority 8 / max retries 3</span>
              </div>
              <strong>running</strong>
            </div>

            <div className="task-card">
              <div>
                <p>send_email</p>
                <span>payload validated</span>
              </div>
              <strong>pending</strong>
            </div>

            <div className="timeline">
              <div>
                <Clock3 size={16} />
                <span>Created</span>
                <time>12:04</time>
              </div>
              <div>
                <Activity size={16} />
                <span>Running</span>
                <time>12:05</time>
              </div>
              <div>
                <ShieldCheck size={16} />
                <span>Auditable</span>
                <time>ready</time>
              </div>
            </div>
          </aside>
        </div>

        <section className="stats-row" aria-label="Platform summary">
          {stats.map((stat) => (
            <div className="stat-item" key={stat.label}>
              <strong>{stat.value}</strong>
              <span>{stat.label}</span>
              <p>{stat.detail}</p>
            </div>
          ))}
        </section>

        <section className="action-grid" aria-label="Primary actions">
          {actions.map((action) => {
            const Icon = action.icon

            return (
              <a
                className={action.primary ? 'action-card featured' : 'action-card'}
                href={action.href}
                key={action.title}
              >
                <span className="action-icon">
                  <Icon size={20} />
                </span>
                <span>
                  <strong>{action.title}</strong>
                  <small>{action.description}</small>
                </span>
                <ArrowRight className="action-arrow" size={18} />
              </a>
            )
          })}
        </section>

        <section className="workflow-section" id="dashboard">
          <div>
            <p className="section-kicker">Next build target</p>
            <h2>One theme, then the real dashboard.</h2>
          </div>

          <div className="workflow-list">
            {workflow.map((item) => {
              const Icon = item.icon

              return (
                <article className="workflow-item" key={item.title}>
                  <Icon size={20} />
                  <div>
                    <h3>{item.title}</h3>
                    <p>{item.text}</p>
                  </div>
                </article>
              )
            })}
          </div>
        </section>
      </section>
    </main>
  )
}

export default App
