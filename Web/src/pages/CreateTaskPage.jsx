import { useMemo, useState } from 'react'
import toast from 'react-hot-toast'
import {
  ArrowLeft,
  Braces,
  CheckCircle2,
  CircleAlert,
  Copy,
  FileJson,
  Gauge,
  Plus,
  RotateCcw,
  Send,
} from 'lucide-react'
import { createTask } from '@/api/tasks'
import TopNav from '@/components/layout/TopNav'

const payloadTemplates = {
  send_email: {
    to: 'team@taskforge.dev',
    subject: 'Worker queue update',
    body: 'The task has been queued successfully.',
  },
  generate_report: {
    report_id: 'weekly-summary',
    format: 'pdf',
    include_charts: true,
  },
  custom: {
    source: 'manual',
    metadata: {
      requested_by: 'operator',
    },
  },
}

function formatJson(value) {
  return JSON.stringify(value, null, 2)
}

function CreateTaskPage() {
  const [taskType, setTaskType] = useState('send_email')
  const [priority, setPriority] = useState(5)
  const [maxRetries, setMaxRetries] = useState(3)
  const [payload, setPayload] = useState(formatJson(payloadTemplates.send_email))
  const [createdTask, setCreatedTask] = useState(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const parsedPayload = useMemo(() => {
    try {
      return {
        data: JSON.parse(payload),
        error: '',
      }
    } catch (error) {
      return {
        data: null,
        error: error.message,
      }
    }
  }, [payload])

  const taskPreview = useMemo(
    () => ({
      type: taskType,
      payload: parsedPayload.data || payload,
      priority: Number(priority),
      max_retries: Number(maxRetries),
    }),
    [maxRetries, parsedPayload.data, payload, priority, taskType],
  )

  const handleTypeChange = (event) => {
    const nextType = event.target.value

    setTaskType(nextType)
    setPayload(formatJson(payloadTemplates[nextType]))
    setCreatedTask(null)
  }

  const handleSubmit = async (event) => {
    event.preventDefault()

    if (parsedPayload.error) {
      toast.error('Fix the payload JSON before submitting.')
      return
    }

    setIsSubmitting(true)

    try {
      const result = await toast.promise(createTask(taskPreview), {
        loading: 'Creating task...',
        success: ({ taskDetails }) => {
          const shortId = taskDetails?.id ? taskDetails.id.slice(0, 8) : 'created'

          return `Task ${shortId} queued.`
        },
        error: (error) => error.message || 'Could not create task.',
      })

      setCreatedTask(result)
    } catch {
      setCreatedTask(null)
    } finally {
      setIsSubmitting(false)
    }
  }

  const resetForm = () => {
    setTaskType('send_email')
    setPriority(5)
    setMaxRetries(3)
    setPayload(formatJson(payloadTemplates.send_email))
    setCreatedTask(null)
  }

  return (
    <section className="app-shell">
      <TopNav />

      <div className="page-heading ">
        <a className="back-link m-5" href="#">
          <ArrowLeft size={16} />
          Home
        </a>
        <div className="eyebrow">
          <Plus size={16} />
          Create task
        </div>
        <h1>Shape a task before it enters the queue.</h1>
        <p>
          Configure the task type, priority, retry policy, and JSON payload.
          Submitting this form sends the request body to the backend create-task endpoint.
        </p>
      </div>

      <div className="create-grid">
        <form className="task-form" onSubmit={handleSubmit}>
          <div className="form-section-header">
            <div>
              <span>Task configuration</span>
              <h2>Queue input</h2>
            </div>
            <FileJson size={22} />
          </div>

          <label className="field-group">
            <span>Task type</span>
            <select value={taskType} onChange={handleTypeChange}>
              <option value="send_email">send_email</option>
              <option value="generate_report">generate_report</option>
              <option value="custom">custom</option>
            </select>
          </label>

          <div className="field-row">
            <label className="field-group">
              <span>Priority</span>
              <input
                min="0"
                max="10"
                type="number"
                value={priority}
                onChange={(event) => setPriority(event.target.value)}
              />
            </label>

            <label className="field-group">
              <span>Max retries</span>
              <input
                min="0"
                max="10"
                type="number"
                value={maxRetries}
                onChange={(event) => setMaxRetries(event.target.value)}
              />
            </label>
          </div>

          <label className="field-group">
            <span>Payload JSON</span>
            <textarea
              spellCheck="false"
              value={payload}
              onChange={(event) => setPayload(event.target.value)}
            />
          </label>

          <div className={parsedPayload.error ? 'json-status error' : 'json-status'}>
            {parsedPayload.error ? <CircleAlert size={17} /> : <CheckCircle2 size={17} />}
            <span>{parsedPayload.error || 'Payload is valid JSON'}</span>
          </div>

          <div className="form-actions">
            <button className="button button-secondary" type="button" onClick={resetForm}>
              Reset
              <RotateCcw size={17} />
            </button>
            <button
              className="button button-primary"
              type="submit"
              disabled={Boolean(parsedPayload.error) || isSubmitting}
            >
              {isSubmitting ? 'Creating...' : 'Submit Task'}
              <Send size={17} />
            </button>
          </div>
        </form>

        <aside className="preview-panel" aria-label="Task request preview">
          <div className="console-header">
            <span>POST /api/task</span>
            <div className="status-pill">
              <span></span>
              draft
            </div>
          </div>

          <div className="preview-summary">
            <div>
              <Gauge size={18} />
              <span>Priority</span>
              <strong>{priority}</strong>
            </div>
            <div>
              <RotateCcw size={18} />
              <span>Retries</span>
              <strong>{maxRetries}</strong>
            </div>
            <div>
              <Braces size={18} />
              <span>Payload</span>
              <strong>{parsedPayload.error ? 'invalid' : 'valid'}</strong>
            </div>
          </div>

          <div className="code-card">
            <div className="code-card-header">
              <span>request body</span>
              <Copy size={16} />
            </div>
            <pre>{formatJson(taskPreview)}</pre>
          </div>

          {createdTask ? (
            <div className="created-task-card">
              <div className="code-card-header">
                <span>created task</span>
                <strong>{createdTask.status}</strong>
              </div>
              <pre>{formatJson(createdTask.taskDetails)}</pre>
            </div>
          ) : null}
        </aside>
      </div>
    </section>
  )
}

export default CreateTaskPage
