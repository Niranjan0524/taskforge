import { useCallback, useEffect, useMemo, useState } from 'react'
import  {
    ReactFlow,
  Background,
  BackgroundVariant,
  Controls,
  MiniMap,
  useEdgesState,
  useNodesState,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import QueueNode from './QueueNode'
import TaskDetailsDrawer from './TaskDetailsDrawer'
import {
  buildQueueNodeData,
  createInitialQueueNodes,
  FLOW_EDGES,
  FLOW_EDGES_RUNNING,
  groupTasksByStatus,
} from './flowUtils'
import './TaskFlow.css'

const nodeTypes = {
  queue: QueueNode,
}

const initialNodes = createInitialQueueNodes()

function TaskFlow({ tasks = []  }) {
  const [selectedTask, setSelectedTask] = useState(null)
  
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(FLOW_EDGES_RUNNING)

  const handleTaskClick = useCallback((task) => {
    setSelectedTask(task)
  }, [])

  const tasksByStatus = useMemo(() => groupTasksByStatus(tasks), [tasks])

  const queueNodeData = useMemo(
    () => buildQueueNodeData(tasksByStatus, handleTaskClick),
    [handleTaskClick, tasksByStatus],
  )

  useEffect(() => {
    setNodes((currentNodes) =>
      currentNodes.map((node) => {
        const nextData = queueNodeData[node.id]

        return nextData
          ? {
              ...node,
              data: nextData,
            }
          : node
      }),
    )
  }, [queueNodeData, setNodes])

  const onConnect = useCallback(
    (connection) => {
      setEdges((eds) => [...eds, connection])
    },
    [setEdges],
  )

  return (
    <section className="taskflow-shell" aria-label="Task queue flow">
      <div className="taskflow-canvas">
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
          minZoom={0.55}
          maxZoom={1.8}
          defaultEdgeOptions={{ type: 'smoothstep' }}
        >
          <Background
            gap={40}
            size={1.25}
            variant={BackgroundVariant.Dots}
            className="task-background"
          />
          <Controls showInteractive={false} />
          <MiniMap
            nodeColor={(node) => node.data?.theme?.accent || '#64748b'}
            maskColor="rgba(37, 38, 40, 0.5)"
            className="task-minimap"
            pannable
            zoomable
          />
        </ReactFlow>
      </div>

      <TaskDetailsDrawer
        open={Boolean(selectedTask)}
        task={selectedTask}
        onClose={() => setSelectedTask(null)}
      />
    </section>
  )
}

export default TaskFlow
