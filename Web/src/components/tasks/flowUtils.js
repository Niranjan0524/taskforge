export const FLOW_STATUSES = ['pending', 'running', 'completed', 'failed', 'cancelled']

export const FLOW_STATUS_META = {
  pending: {
    label: 'Pending',
    accent: '#f6c94c',
    tint: 'rgba(246, 201, 76, 0.14)',
  },
  running: {
    label: 'Running',
    accent: '#4ea0ff',
    tint: 'rgba(78, 160, 255, 0.14)',
  },
  completed: {
    label: 'Completed',
    accent: '#4fd08b',
    tint: 'rgba(79, 208, 139, 0.14)',
  },
  failed: {
    label: 'Failed',
    accent: '#ff6b6b',
    tint: 'rgba(255, 107, 107, 0.14)',
  },
  cancelled: {
    label: 'Cancelled',
    accent: '#98a2b3',
    tint: 'rgba(152, 162, 179, 0.14)',
  },
}

export const FLOW_NODE_POSITIONS = {
  pending: { x: 0, y: 0 },
  running: { x: 370, y: 0 },
  completed: { x: 740, y: -210 },
  failed: { x: 740, y: 0 },
  cancelled: { x: 740, y: 210 },
}

export const FLOW_EDGES = [
  {
    id: 'pending-running',
    source: 'pending',
    target: 'running',
    animated: true,
  },
  {
    id: 'running-completed',
    source: 'running',
    target: 'completed',
    animated: true,
  },
  {
    id: 'running-failed',
    source: 'running',
    target: 'failed',
    animated: true,
  },
  {
    id: 'running-cancelled',
    source: 'running',
    target: 'cancelled',
    animated: true,
  },
]

export function normalizeTaskStatus(status) {
  const normalized = String(status || 'pending').toLowerCase()

  return FLOW_STATUSES.includes(normalized) ? normalized : 'pending'
}

export function shortTaskId(taskId) {
  if (!taskId) {
    return 'unknown'
  }

  return taskId.length > 8 ? taskId.slice(0, 8) : taskId
}

export function groupTasksByStatus(tasks = []) {
  return tasks.reduce(
    (accumulator, task) => {
      const status = normalizeTaskStatus(task?.status)

      accumulator[status].push(task)

      return accumulator
    },
    {
      pending: [],
      running: [],
      completed: [],
      failed: [],
      cancelled: [],
    },
  )
}

export function buildQueueNodeData(tasksByStatus, onTaskClick) {
  return FLOW_STATUSES.reduce((accumulator, status) => {
    const tasks = tasksByStatus[status] || []

    accumulator[status] = {
      title: FLOW_STATUS_META[status].label,
      status,
      tasks,
      theme: FLOW_STATUS_META[status],
      onTaskClick,
    }

    return accumulator
  }, {})
}

export function createInitialQueueNodes() {
  return FLOW_STATUSES.map((status) => ({
    id: status,
    type: 'queue',
    data: {
      title: FLOW_STATUS_META[status].label,
      status,
      tasks: [],
      theme: FLOW_STATUS_META[status],
      onTaskClick: null,
    },
    position: FLOW_NODE_POSITIONS[status],
  }))
}