# TaskForge

TaskForge is a distributed background task processing system built with **Go** and **Redis**. It enables applications to execute long-running tasks asynchronously using a worker pool, providing real-time task monitoring through WebSockets.

Instead of executing time-consuming operations synchronously, TaskForge pushes tasks into a Redis-backed queue where worker goroutines process them independently, allowing APIs to respond immediately while work continues in the background.

---

## Features

* Asynchronous background task processing
* Processing tasks based in priority
* Redis-based task queue
* Concurrent worker pool using Goroutines
* Real-time task status updates via WebSockets
* Task cancellation support
* Random failure simulation for realistic task execution
* JWT-protected REST APIs
* Redis Pub/Sub integration for live updates
* Retry mechanism with exponential backoff

---

## Architecture

```
                 Client
                    │
        REST API / WebSocket
                    │
             Gin HTTP Server
                    │
             Task Management
                    │
                 Redis
      ┌─────────────┴─────────────┐
      │                           │
 Task Metadata              Priority Queue
      │                           │
      └─────────────┬─────────────┘
                    │
             Worker Pool (Go)
                    │
          Task Execution Engine
                    │
        WebSocket Status Updates
```

---

## Supported Task Types

Current simulated task handlers:

* Send Email
* Generate Report
* Resize Image

Each task is executed asynchronously by workers and transitions through its lifecycle independently.

---

## Task Lifecycle

```
Pending
   │
   ▼
Running
   │
   ├────────► Cancelled
   │
   ├────────► Failed ------->Dead
   │
   ▼
Completed
```

---

## REST APIs

| Method | Endpoint                | Description                         |
| ------ | ----------------------- | ----------------------------------- |
| POST   | `/api/task`             | Create a new task                   |
| GET    | `/api/task/:id`         | Get task details                    |
| GET    | `/api/task/:id/status`  | Get task status                     |
| GET    | `/api/tasks`            | Get all tasks                       |
| DELETE | `/api/tasks/:id`        | Delete a task                       |
| DELETE | `/api/tasks/cancel/:id` | Cancel a running or pending task    |
| GET    | `/ws`                   | WebSocket endpoint for live updates |

---

## Tech Stack

### Backend

* Go
* Gin
* Redis
* Gorilla WebSocket

### Concurrency

* Goroutines
* Channels
* Worker Pool

### Communication

* REST APIs
* WebSockets
* Redis Pub/Sub

---

## Project Structure

```
server/
│
├── cmd/
│   ├── api/
│   └── worker/
│
├── internals/
│   ├── handlers/
│   ├── Storage/
│   ├── worker/
│   ├── models/
│   ├── queue/
│   └── webSockets/
│
└── main.go
```

---

## Running the Project

### Clone the repository

```bash
git clone <repository-url>
cd taskforge
```

### Start Redis

```bash
docker run -d \
-p 6379:6379 \
--name redis \
redis:latest
```

### Configure Environment

Create a `.env` file.

```
REDIS_ADDR=localhost:6379
ORIGIN_URL=http://localhost:5173
JWT_SECRET=your_secret
```

### Run API Server

```bash
go run main.go
```

---

## WebSocket Events

TaskForge broadcasts task updates in real time.

Example message:

```json
{
  "taskId": "8c44f0c8",
  "status": "running"
}
```

Example statuses:

* pending
* running
* completed
* failed
* cancelled
* dead

---

## Future Improvements


* Scheduled and delayed tasks
* Task prioritization
* Prometheus metrics
* Grafana dashboards
* Docker Compose deployment
* Rate limiting

---

## Learning Outcomes

This project demonstrates:

* Concurrent programming with Goroutines
* Producer–Consumer architecture
* Redis data structures and Pub/Sub
* Asynchronous task execution
* Real-time communication using WebSockets
* REST API design
* Background job processing
* Distributed system fundamentals
