package main

import (
	"api/internal/services"
	"api/pkg/utils"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	utils.Migrate()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"nodename": os.Getenv("nodename"),
		})
	})

	r.POST("/register", services.Register)
	r.POST("/authenticate", services.Authenticate)
	r.GET("/tasks", services.GetTaskForUser)
	r.POST("/tasks", services.CreateTaskForUser)
	r.DELETE("/tasks/:id", services.DeleteTaskForUser)

	r.Run(":8080")
}
