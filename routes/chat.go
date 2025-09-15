package routes

import (
	"github.com/gin-gonic/gin"
	"temp/controllers"
	"temp/middleware"
)

func ChatRoutes(r *gin.Engine) {
	chat := r.Group("/chats")
	{
		chat.POST("/message", middleware.AuthMiddleware(), controllers.SendMessage)
		chat.GET("/:chatId/messages", middleware.AuthMiddleware(), controllers.GetMessages)

	}
}
