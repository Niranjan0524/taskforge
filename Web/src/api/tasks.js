const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || ''

async function readResponse(response) {
  const contentType = response.headers.get('content-type') || ''

  if (contentType.includes('application/json')) {
    return response.json()
  }

  return response.text()
}

export async function createTask(task) {
  const response = await fetch(`${API_BASE_URL}/api/task`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(task),
  })

  const data = await readResponse(response)

  if (!response.ok) {
    const message = data?.error || data || 'Failed to create task'
    throw new Error(message)
  }

  return {
    status: response.status,
    taskDetails: data,
  }
}

export async function getTasks() {
  const response = await fetch(`${API_BASE_URL}/api/tasks`)
  const data = await readResponse(response)

  if (!response.ok) {
    const message = data?.error || data || 'Failed to fetch tasks'
    throw new Error(message)
  }

  return Array.isArray(data) ? data : []
}

export async function getTask(taskId) {
  const response = await fetch(`${API_BASE_URL}/api/task/${taskId}`)
  const data = await readResponse(response)

  if (!response.ok) {
    const message = data?.error || data || 'Failed to fetch task'
    throw new Error(message)
  }

  return data
}

export async function getTaskStatus(taskId) {
  const response = await fetch(`${API_BASE_URL}/api/task/${taskId}/status`)
  const data = await readResponse(response)

  if (!response.ok) {
    const message = data?.error || data || 'Failed to refresh task status'
    throw new Error(message)
  }

  return data
}

export async function deleteTask(taskId) {
  const response = await fetch(`${API_BASE_URL}/api/tasks/${taskId}`, {
    method: 'DELETE',
  })

  const data = await readResponse(response)

  if (!response.ok) {
    const message = data?.error || data || 'Failed to delete task'
    throw new Error(message)
  }

  return data
}
