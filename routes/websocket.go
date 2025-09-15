package routes

import (
	"github.com/gin-gonic/gin"
	"temp/websocket"
)

func WebSocketRoutes(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		// Hand off to your existing websocket handler
		websocket.HandleWebSocket(c.Writer, c.Request)
	})
}
