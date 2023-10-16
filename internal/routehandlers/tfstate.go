package routehandlers

import (
	"github.com/arpanrec/secureserver/internal/tfstate"
	"github.com/gin-gonic/gin"
	"io"
)

func TfStateHandler() gin.HandlerFunc {
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
		query := c.Request.URL.Query()
		s, d := tfstate.TerraformStateHandler(string(body), rMethod, locationPath, query)
		c.Data(s, "application/json", []byte(d))
	}
}
