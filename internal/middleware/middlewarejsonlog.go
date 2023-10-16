package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			logMap := make(map[string]interface{})

			logMap["status_code"] = params.StatusCode
			logMap["path"] = params.Path
			logMap["method"] = params.Method
			logMap["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			logMap["remote_addr"] = params.ClientIP
			logMap["response_time"] = params.Latency.String()

			s, err := json.Marshal(logMap)
			if err != nil {
				log.Println("Error while marshalling logMap" + err.Error())
				return "Error while marshalling logMap" + err.Error()
			}
			return string(s) + "\n"
		},
	)
}
