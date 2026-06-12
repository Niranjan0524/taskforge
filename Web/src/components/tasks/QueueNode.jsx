import { Handle, Position } from '@xyflow/react'
import './QueueNode.css'

function QueueNode({ data }) {
  const title = data?.title || 'Status'
  const tasks = Array.isArray(data?.tasks) ? data.tasks : []
  const taskCount = tasks.length

  return (
    <div className="queue-node">
      <Handle type="target" position={Position.Left} className="queue-handle" />
      <div className="queue-node-body">
        <span className="queue-node-title">{title}</span>
        <strong className="queue-node-count">
          {taskCount} {taskCount === 1 ? 'Task' : 'Tasks'}
        </strong>
      </div>
      <Handle type="source" position={Position.Right} className="queue-handle" />
    </div>
  )
}

export default QueueNode