package ginhosting

import (
	"github.com/arpanrec/secureserver/internal/pki"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func PkiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, errReadAll := io.ReadAll(c.Request.Body)
		if errReadAll != nil {
			c.JSON(500, gin.H{
				"error": errReadAll.Error(),
			})
			return
		}
		locationPath := c.GetString("locationPath")

		r, err := pki.GetCert(&locationPath, &body)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(http.StatusCreated, "application/json", []byte(r))
	}
}
