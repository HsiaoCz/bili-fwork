package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/user", Hello)
	r.Run(":9092")
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Hello",
	})
}
