package ginhosting

import (
	"github.com/arpanrec/secretsquirrel/internal/api"
	"github.com/gin-gonic/gin"
	"io"
)

func keyValueHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, errReadAll := io.ReadAll(c.Request.Body)
		if errReadAll != nil {
			c.JSON(500, gin.H{
				"error": errReadAll.Error(),
			})
			return
		}
		rMethod := c.Request.Method
		key := c.Request.URL.Path[7:]
		kvData, err := api.KeyValue(&body, &rMethod, &key)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if kvData == nil {
			c.Data(200, "application/json", nil)
			return
		}
		c.JSON(200, kvData)
		return
	}
}
