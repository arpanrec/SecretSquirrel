package middleware

import (
	"encoding/json"
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
			ginCtxErr := c.AbortWithError(http.StatusUnauthorized, err)
			log.Println("Error in AuthMiddleWare: ", ginCtxErr)
		}
		c.Set("username", username)
		return
	}
}

func NameSpaceMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		urlPath := c.Request.URL.Path
		locationPath := username + urlPath[7:]
		c.Set("locationPath", locationPath)
		return
	}
}

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
