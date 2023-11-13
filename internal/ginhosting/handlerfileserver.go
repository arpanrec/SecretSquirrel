package ginhosting

import (
	"io"

	"github.com/arpanrec/secretsquirrel/internal/fileserver"
	"github.com/gin-gonic/gin"
)

func fileServerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, errReadAll := io.ReadAll(c.Request.Body)
		if errReadAll != nil {
			c.JSON(500, gin.H{
				"error": errReadAll.Error(),
			})
			return
		}
		rMethod := c.Request.Method
		locationPath := c.GetString("locationPath")
		s, d := fileserver.ReadWriteFilesFromURL(string(body), rMethod, locationPath)
		c.Data(s, "text/plain", []byte(d))
	}
}
