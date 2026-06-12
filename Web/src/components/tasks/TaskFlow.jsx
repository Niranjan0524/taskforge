import { useCallback, useEffect, useMemo, useState } from 'react'
import ReactFlow, {
  Background,
  BackgroundVariant,
  Controls,
  MiniMap,
  useEdgesState,
  useNodesState,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { getTasks } from '@/api/tasks'
import QueueNode from './QueueNode'
import './TaskFlow.css'

const statusOrder = ['pending', 'running', 'completed', 'failed', 'cancelled']

const statusLabels = {
  pending: 'Pending',
  running: 'Running',
  completed: 'Completed',
  failed: 'Failed',
  cancelled: 'Cancelled',
}

const nodeTypes = {
  queue: QueueNode,
}

const initialNodes = [
  {
    id: 'pending',
    data: { title: 'Pending', tasks: [] },
    position: { x: 0, y: 0 },
    type: 'queue',
  },
  {
    id: 'running',
    data: { title: 'Running', tasks: [] },
    position: { x: 320, y: 0 },
    type: 'queue',
  },
  {
    id: 'completed',
    data: { title: 'Completed', tasks: [] },
    position: { x: 640, y: -180 },
    type: 'queue',
  },
  {
    id: 'failed',
    data: { title: 'Failed', tasks: [] },
    position: { x: 640, y: 0 },
    type: 'queue',
  },
  {
    id: 'cancelled',
    data: { title: 'Cancelled', tasks: [] },
    position: { x: 640, y: 180 },
    type: 'queue',
  },
]

const initialEdges = [
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

function groupTasksByStatus(tasks) {
  return tasks.reduce(
    (accumulator, task) => {
      const normalizedStatus = String(task?.status || 'pending').toLowerCase()

      if (statusOrder.includes(normalizedStatus)) {
        accumulator[normalizedStatus].push(task)
      }

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

function buildQueueNodes(tasksByStatus) {
  return statusOrder.map((status, index) => ({
    id: status,
    type: 'queue',
    data: {
      title: statusLabels[status],
      tasks: tasksByStatus[status],
    },
    position: initialNodes[index].position,
  }))
}

function TaskFlow() {
  const [tasks, setTasks] = useState([])
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)

  useEffect(() => {
    let isMounted = true

    async function loadTasks() {
      try {
        const taskList = await getTasks()

        if (isMounted) {
          setTasks(taskList)
        }
      } catch {
        if (isMounted) {
          setTasks([])
        }
      }
    }

    loadTasks()

    const socket = new WebSocket('ws://localhost:8080/ws')

    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)

        if (!message?.taskId || !message?.status) {
          return
        }

        setTasks((currentTasks) =>
          currentTasks.map((task) =>
            task.id === message.taskId
              ? {
                  ...task,
                  status: message.status,
                }
              : task,
          ),
        )
      } catch {
        return
      }
    }

    return () => {
      isMounted = false
      socket.close()
    }
  }, [])

  const tasksByStatus = useMemo(() => groupTasksByStatus(tasks), [tasks])

  const queueNodes = useMemo(
    () => buildQueueNodes(tasksByStatus),
    [tasksByStatus],
  )

  useEffect(() => {
    setNodes(queueNodes)
  }, [queueNodes, setNodes])

  const onConnect = useCallback(
    (connection) => {
      setEdges((eds) => [...eds, connection])
    },
    [setEdges],
  )

  return (
    <div className="taskflow-container">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        fitView
        panOnDrag
        zoomOnScroll
        zoomOnPinch
        minZoom={0.6}
        maxZoom={1.8}
      >
        <Background
          color="#233044"
          gap={18}
          size={1}
          variant={BackgroundVariant.Dots}
          className="task-background"
        />
        <Controls showInteractive={false} />
        <MiniMap
          nodeColor="#2b3b57"
          maskColor="rgba(5, 10, 20, 0.55)"
          className="task-minimap"
          pannable
          zoomable
        />
      </ReactFlow>
    </div>
  )
}

export default TaskFlow
