package webSockets

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == os.Getenv("ORIGIN_URL")
	},
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		conn.WriteMessage(messageType, msg)
	}
}
