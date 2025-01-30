package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	answeringNode := os.Getenv("nodename")

	html := "<html><body>"
	html += "<h1>Hello, World!</h1>"
	html += answeringNode
	html += "</body></html>"

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})
	r.Run(":8080")
}
