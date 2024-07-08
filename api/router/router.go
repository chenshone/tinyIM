package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}
