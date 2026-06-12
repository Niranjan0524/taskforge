import { memo, useEffect, useMemo, useState } from 'react'
import { Handle, Position } from '@xyflow/react'
import { shortTaskId } from './flowUtils'
import './QueueNode.css'

function buildTaskSignature(tasks) {
  return tasks
    .map((task) => [task?.id, task?.status, task?.priority, task?.retry_count, task?.type].join(':'))
    .join('|')
}

function QueueNode({ data, selected }) {
  const title = data?.title || 'Status'
  const tasks = Array.isArray(data?.tasks) ? data.tasks : []
  const visibleTasks = useMemo(() => tasks.slice(0, 5), [tasks])
  const moreCount = Math.max(tasks.length - visibleTasks.length, 0)
  const taskCount = tasks.length
  const theme = data?.theme || {}
  const onTaskClick = data?.onTaskClick
  const [pulse, setPulse] = useState(false)

  const taskSignature = useMemo(() => buildTaskSignature(tasks), [tasks])

  useEffect(() => {
    setPulse(true)

    const timeout = window.setTimeout(() => {
      setPulse(false)
    }, 260)

    return () => window.clearTimeout(timeout)
  }, [taskSignature])

  return (
    <div
      className={`queue-node queue-node--${data?.status || 'pending'}${selected ? ' is-selected' : ''}${pulse ? ' is-updated' : ''}`}
      style={{ '--status-accent': theme.accent, '--status-tint': theme.tint }}
    >
      <Handle type="target" position={Position.Left} className="queue-handle" />

      <div className="queue-node-body">
        <div className="queue-node-header">
          <span className="queue-node-title">{title}</span>
          <strong className="queue-node-count">
            {taskCount} {taskCount === 1 ? 'Task' : 'Tasks'}
          </strong>
        </div>

        <div className="queue-task-list">
          {visibleTasks.length ? (
            visibleTasks.map((task) => (
              <button
                key={task.id}
                type="button"
                className="queue-task-card"
                onClick={() => onTaskClick?.(task)}
              >
                <div className="queue-task-row">
                  <span className="queue-task-id">{shortTaskId(task.id)}</span>
                  <span className="queue-task-priority">P{task.priority ?? 0}</span>
                </div>
                <div className="queue-task-row queue-task-row--subtle">
                  <span className="queue-task-type">{task.type || 'unknown'}</span>
                  <span className="queue-task-retry">
                    Retry {task.retry_count ?? 0}/{task.max_retries ?? 0}
                  </span>
                </div>
              </button>
            ))
          ) : (
            <span className="queue-node-empty">No tasks</span>
          )}

          {moreCount > 0 ? <span className="queue-node-more">+{moreCount} more</span> : null}
        </div>
      </div>

      <Handle type="source" position={Position.Right} className="queue-handle" />
    </div>
  )
}

export default memo(QueueNode)