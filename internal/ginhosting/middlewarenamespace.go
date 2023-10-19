package ginhosting

import (
	"log"

	"github.com/gin-gonic/gin"
)

func nameSpaceMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		urlPath := c.Request.URL.Path
		locationPath := username + urlPath[7:]
		c.Set("locationPath", locationPath)
		log.Println("nameSpaceMiddleWare: Namespace is " + locationPath)
	}
}
