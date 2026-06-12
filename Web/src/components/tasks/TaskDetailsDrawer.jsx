import { memo } from 'react'
import { X } from 'lucide-react'
import './TaskDetailsDrawer.css'

function TaskDetailsDrawer({ open, task, onClose }) {
  if (!open || !task) {
    return null
  }

  return (
    <aside className="task-details-drawer" aria-label="Task details panel">
      <div className="task-details-header">
        <div>
          <span className="task-details-eyebrow">Selected task</span>
          <h3>{task.type || 'Task details'}</h3>
        </div>

        <button
          className="task-details-close"
          type="button"
          aria-label="Close task details"
          onClick={onClose}
        >
          <X size={18} />
        </button>
      </div>

      <div className="task-details-meta">
        <div>
          <span>Task ID</span>
          <strong>{task.id}</strong>
        </div>
        <div>
          <span>Status</span>
          <strong>{task.status}</strong>
        </div>
        <div>
          <span>Priority</span>
          <strong>{task.priority ?? 0}</strong>
        </div>
        <div>
          <span>Retries</span>
          <strong>{task.retry_count ?? 0} / {task.max_retries ?? 0}</strong>
        </div>
      </div>

      <div className="task-details-json">
        <pre>{JSON.stringify(task, null, 2)}</pre>
      </div>
    </aside>
  )
}

export default memo(TaskDetailsDrawer)