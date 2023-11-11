package ginhosting

import (
	"log"
	"net/http"

	"github.com/arpanrec/secretsquirrel/internal/auth"
	"github.com/gin-gonic/gin"
)

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("authMiddleWare")
		authHeader := c.GetHeader("Authorization")
		username, err := auth.GetUserDetails(authHeader)
		if err != nil {
			// Abort and set 401 and set the error in body
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		c.Set("username", username)
	}
}
