import { BarChart3, ListChecks } from 'lucide-react'
import TopNav from '@/components/layout/TopNav'
import TaskList from '@/components/tasks/TaskList'

function DashboardPage() {
  return (
    <section className="app-shell">
      <TopNav />

      <div className="page-heading dashboard-heading">
        <div className="eyebrow">
          <BarChart3 size={16} />
          Dashboard
        </div>
        <h1>Track every task moving through the queue.</h1>
        <p>
          Review task type, status, priority, retry progress, creation time,
          and the core actions needed while the worker pool is running.
        </p>
      </div>

      <section className="dashboard-panel">
        <div className="form-section-header">
          <div>
            <span>Task inventory</span>
            <h2>Task list</h2>
          </div>
          <ListChecks size={22} />
        </div>

        <TaskList />
      </section>
    </section>
  )
}

export default DashboardPage
