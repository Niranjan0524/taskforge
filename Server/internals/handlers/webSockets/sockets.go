package webSockets

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
	"github.com/Niranjan0524/taskforge/server/internals/config"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const TaskStatusChannel = "tasks:status"

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return config.OriginAllowed(r.Header.Get("Origin"))
	},
}

type Hub struct {
	Clients          map[*websocket.Conn]bool
	Broadcast        chan []byte
	mu               sync.RWMutex
	lastWorkerStatus string
}

var WsHub = &Hub{
	Clients:   make(map[*websocket.Conn]bool),
	Broadcast: make(chan []byte, 64),
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func() {
		WsHub.mu.Lock()
		delete(WsHub.Clients, conn)
		WsHub.mu.Unlock()
		conn.Close()
	}()

	WsHub.mu.Lock()
	WsHub.Clients[conn] = true
	status := WsHub.lastWorkerStatus
	WsHub.mu.Unlock()

	if status != "" {
		if data, err := MarshalWorkerStatus(status); err == nil {
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			break
		}
	}
}

func (h *Hub) Run() {
	for msg := range h.Broadcast {
		h.mu.RLock()
		clients := make([]*websocket.Conn, 0, len(h.Clients))
		for conn := range h.Clients {
			clients = append(clients, conn)
		}
		h.mu.RUnlock()

		for _, conn := range clients {
			err := conn.WriteMessage(
				websocket.TextMessage,
				msg,
			)
			if err != nil {
				conn.Close()
				h.mu.Lock()
				delete(h.Clients, conn)
				h.mu.Unlock()
			}
		}
	}
}

type wsMeassage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
type TaskStatusMessage struct {
	TaskID string       `json:"taskId"`
	Task   storage.Task `json:"task"`
}

type WorkerStatusMessage struct {
	Status string `json:"status"`
}

func MarshalTaskStatus(taskID string, task storage.Task) ([]byte, error) {
	msg := TaskStatusMessage{
		TaskID: strings.TrimPrefix(taskID, "task:"),
		Task:   task,
	}

	wsMsg := wsMeassage{
		Type: "taskUpdate",
		Data: msg,
	}
	return json.Marshal(wsMsg)
}

func MarshalWorkerStatus(status string) ([]byte, error) {
	data := WorkerStatusMessage{
		Status: status,
	}

	wsMsg := wsMeassage{
		Type: "workerStatus",
		Data: data,
	}
	return json.Marshal(wsMsg)
}

func BroadcastRaw(data []byte) {

	select {
	case WsHub.Broadcast <- data:
	default:
		log.Println("websocket broadcast channel full; dropping task status update")
	}
}

func PublishWorkerStatus(ctx context.Context, client *redis.Client, status string) error {
	data, err := MarshalWorkerStatus(status)
	if err != nil {
		return err
	}
	return client.Publish(ctx, TaskStatusChannel, data).Err()
}

func BroadcastTaskStatus(
	taskID string,
	task storage.Task,
) {

	data, err := MarshalTaskStatus(taskID, task)
	if err != nil {
		return
	}

	BroadcastRaw(data)
}

func BroadcastWorkerStatus(
	status string,
) {

	data, err := MarshalWorkerStatus(status)
	if err != nil {
		return
	}
	WsHub.mu.Lock()
	WsHub.lastWorkerStatus = status
	WsHub.mu.Unlock()
	BroadcastRaw(data)
}

func StartTaskStatusSubscriber(ctx context.Context, client *redis.Client) {
	pubsub := client.Subscribe(ctx, TaskStatusChannel)
	defer pubsub.Close()

	channel := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-channel:
			if !ok {
				return
			}
			payload := []byte(msg.Payload)

			// cache the latest worker status in THIS process's hub so any
			// client that connects later still gets the current state
			var envelope wsMeassage
			if err := json.Unmarshal(payload, &envelope); err == nil && envelope.Type == "workerStatus" {
				if dataMap, ok := envelope.Data.(map[string]interface{}); ok {
					if status, ok := dataMap["status"].(string); ok {
						WsHub.mu.Lock()
						WsHub.lastWorkerStatus = status
						WsHub.mu.Unlock()
					}
				}
			}

			BroadcastRaw(payload)
		}
	}
}
