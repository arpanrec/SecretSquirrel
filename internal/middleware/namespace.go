package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
)

func NameSpaceMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		urlPath := c.Request.URL.Path
		locationPath := username + urlPath[7:]
		c.Set("locationPath", locationPath)
		log.Println("NameSpaceMiddleWare: Namespace is " + locationPath)
		return
	}
}
