package middleware

import (
	"github.com/arpanrec/secureserver/internal/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("AuthMiddleWare")
		authHeader := c.GetHeader("Authorization")
		username, err := auth.GetUserDetails(authHeader)
		if err != nil {
			// Abort and set 401 and set the error in body
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		c.Set("username", username)
		return
	}
}
