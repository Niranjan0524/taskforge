import { useEffect, useMemo, useState } from 'react'
import toast from 'react-hot-toast'
import {
  AlertCircle,
  Eye,
  Loader2,
  RefreshCcw,
  RotateCcw,
  Trash2,
} from 'lucide-react'
import {
  deleteTask,
  getTask,
  getTasks,
  getTaskStatus,
} from '@/api/tasks'

const statusLabels = {
  pending: 'Pending',
  running: 'Running',
  completed: 'Completed',
  failed: 'Failed',
}

function formatDate(value) {
  if (!value) {
    return 'not set'
  }

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return 'invalid date'
  }

  return new Intl.DateTimeFormat('en', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

function shortId(id) {
  return id ? id.slice(0, 8) : 'unknown'
}

function TaskList() {
  const [tasks, setTasks] = useState([])
  const [selectedTask, setSelectedTask] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [activeTaskId, setActiveTaskId] = useState('')

  const taskCount = tasks.length

  const sortedTasks = useMemo(
    () =>
      [...tasks].sort((firstTask, secondTask) => {
        const firstDate = new Date(firstTask.created_at || 0).getTime()
        const secondDate = new Date(secondTask.created_at || 0).getTime()

        return secondDate - firstDate
      }),
    [tasks],
  )

  const loadTasks = async () => {
    setLoading(true)
    setError('')

    try {
      const taskList = await getTasks()

      setTasks(taskList)
    } catch (fetchError) {
      setError(fetchError.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    let isMounted = true

    async function loadInitialTasks() {
      try {
        const taskList = await getTasks()

        if (isMounted) {
          setTasks(taskList)
        }
      } catch (fetchError) {
        if (isMounted) {
          setError(fetchError.message)
        }
      } finally {
        if (isMounted) {
          setLoading(false)
        }
      }
    }

    loadInitialTasks()

    const ws=new WebSocket("ws://localhost:8080/ws")

    ws.onopen=()=>{
      console.log("Connected");
    }

    ws.onmessage=(event)=>{
      const updatedData = JSON.parse(event.data);  // ← Parse the string to object
      console.log(updatedData);

      updateTaskInList(updatedData.taskId, {status: updatedData.status})
      setSelectedTask(null)
      toast.success(`Task ${shortId(taskId)} is ${result.status}.`,{
        position: "bottom-center"
      })
    }

    ws.onclose = () => {
        console.log("Disconnected");
    };

    return () => {
      isMounted = false
      ws.close()
    }
  }, [])



  const updateTaskInList = (taskId, updates) => {
    setTasks((currentTasks) =>
      currentTasks.map((task) =>
        task.id === taskId
          ? {
              ...task,
              ...updates,
            }
          : task,
      ),
    )
  }

  const handleViewTask = async (taskId) => {
    setActiveTaskId(taskId)

    try {
      const task = await getTask(taskId)

      setSelectedTask(task)
      updateTaskInList(taskId, task)
    } catch (viewError) {
      toast.error(viewError.message)
    } finally {
      setActiveTaskId('')
    }
  }

  const handleRefreshStatus = async (taskId) => {
    setActiveTaskId(taskId)

    try {
      const result = await getTaskStatus(taskId)
      updateTaskInList(taskId, { status: result.status })
      setSelectedTask((currentTask) =>
        currentTask?.id === taskId
          ? {
              ...currentTask,
              status: result.status,
            }
          : currentTask,
      )
      toast.success(`Task ${shortId(taskId)} is ${result.status}.`)
    } catch (statusError) {
      toast.error(statusError.message)
    } finally {
      setActiveTaskId('')
    }
  }

  const handleDeleteTask = async (taskId) => {
    const shouldDelete = window.confirm(`Delete task ${shortId(taskId)}?`)

    if (!shouldDelete) {
      return
    }

    setActiveTaskId(taskId)

    try {
      await toast.promise(deleteTask(taskId), {
        loading: 'Deleting task...',
        success: `Task ${shortId(taskId)} deleted.`,
        error: (deleteError) => deleteError.message || 'Could not delete task.',
      })

      setTasks((currentTasks) => currentTasks.filter((task) => task.id !== taskId))
      setSelectedTask((currentTask) => (currentTask?.id === taskId ? null : currentTask))
    } catch {
      return
    } finally {
      setActiveTaskId('')
    }
  }

  if (loading) {
    return (
      <div className="state-panel">
        <Loader2 className="spin-icon" size={22} />
        <strong>Loading tasks</strong>
        <span>Fetching the latest queue state.</span>
      </div>
    )
  }

  if (error) {
    return (
      <div className="state-panel error-state">
        <AlertCircle size={22} />
        <strong>Could not load tasks</strong>
        <span>{error}</span>
        <button className="button button-secondary compact-button" type="button" onClick={loadTasks}>
          Try Again
          <RefreshCcw size={16} />
        </button>
      </div>
    )
  }

  if (!taskCount) {
    return (
      <div className="state-panel">
        <ListEmptyIcon />
        <strong>No tasks yet</strong>
        <span>Create a task to see it appear in this dashboard.</span>
        <a className="button button-primary compact-button" href="#create-task">
          Create Task
        </a>
      </div>
    )
  }

  
  return (
    <div className="task-list-layout">
      <div className="table-toolbar">
        <span>{taskCount} tasks</span>
        <button className="icon-text-button" type="button" onClick={loadTasks}>
          <RefreshCcw size={16} />
          Refresh List
        </button>
      </div>

      <div className="task-table-wrap">
        <table className="task-table">
          <thead>
            <tr>
              <th>Task</th>
              <th>Status</th>
              <th>Priority</th>
              <th>Retries</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {sortedTasks.map((task) => {
              const isActive = activeTaskId === task.id

              return (
                <tr key={task.id}>
                  <td>
                    <strong>{task.type || 'unknown'}</strong>
                    <span>{shortId(task.id)}</span>
                  </td>
                  <td>
                    <span className={`status-badge ${task.status || 'pending'}`}>
                      {statusLabels[task.status] || task.status || 'Pending'}
                    </span>
                  </td>
                  <td>{task.priority ?? 0}</td>
                  <td>
                    {task.retry_count ?? 0} / {task.max_retries ?? 0}
                  </td>
                  <td>{formatDate(task.created_at)}</td>
                  <td>
                    <div className="table-actions">
                      <button
                        aria-label="View task"
                        className="icon-button"
                        disabled={isActive}
                        type="button"
                        onClick={() => handleViewTask(task.id)}
                      >
                        <Eye size={16} />
                      </button>
                      <button
                        aria-label="Refresh task status"
                        className="icon-button"
                        disabled={isActive}
                        type="button"
                        onClick={() => handleRefreshStatus(task.id)}
                      >
                        {isActive ? <Loader2 className="spin-icon" size={16} /> : <RotateCcw size={16} />}
                      </button>
                      <button
                        aria-label="Delete task"
                        className="icon-button danger"
                        disabled={isActive}
                        type="button"
                        onClick={() => handleDeleteTask(task.id)}
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>

      {selectedTask ? (
        <aside className="task-detail-panel">
          <div className="code-card-header">
            <span>selected task</span>
            <strong>{shortId(selectedTask.id)}</strong>
          </div>
          <pre>{JSON.stringify(selectedTask, null, 2)}</pre>
        </aside>
      ) : null}
    </div>
  )
}

function ListEmptyIcon() {
  return (
    <span className="empty-icon" aria-hidden="true">
      <RefreshCcw size={22} />
    </span>
  )
}

export default TaskList
