package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"temp/config"
	"temp/middleware"
	"temp/routes"
	"github.com/gin-contrib/cors"
	"os"
	"github.com/joho/godotenv"
	// "temp/helpers"
)

func main() {
	godotenv.Load()
	// fmt.Println("Calling ConnectDB()...") // Add this
	config.ConnectDB()
	// fmt.Println("Returned from ConnectDB()")
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "PONG"})
	})
	r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Skill Barter API is live!",
    })
})


	r.GET("/protected", middleware.AuthMiddleware(), func(c *gin.Context) {
		email := c.MustGet("email").(string)
		c.JSON(200, gin.H{
			"message": "Welcome to protected route!",
			"email":   email,
		})
	})
	routes.AuthRoutes(r)
	
	port := os.Getenv("PORT")
if port == "" {
    port = "8000" 
}
r.Run(":" + port)


}
