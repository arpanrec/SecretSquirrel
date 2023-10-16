package routehandlers

import (
	"encoding/json"
	"github.com/arpanrec/secureserver/internal/pki"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
)

type pkiRequest struct {
	DnsNames []string `json:"dns_names"`
}

type pkiResponse struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

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
		pkiRequestJson := pkiRequest{}

		err := json.Unmarshal(body, &pkiRequestJson)
		if err != nil {
			c.JSON(500, gin.H{
				"error": errReadAll.Error(),
			})
			return
		}
		log.Println("pkiRequestJson: ", pkiRequestJson)
		var pkiResponseJson pkiResponse
		if strings.HasSuffix(locationPath, "clientcert") {
			cert, k, e := pki.GetClientCert(pkiRequestJson.DnsNames)
			if e != nil {
				c.JSON(500, gin.H{
					"error": e.Error(),
				})
				return
			}
			pkiResponseJson.Cert = cert
			pkiResponseJson.Key = k
		} else if strings.HasSuffix(locationPath, "servercert") {
			cert, k, e := pki.GetServerCert(pkiRequestJson.DnsNames)
			if e != nil {
				c.JSON(500, gin.H{
					"error": e.Error(),
				})
				return
			}
			pkiResponseJson.Cert = cert
			pkiResponseJson.Key = k
		} else {
			c.JSON(500, gin.H{
				"error": "Invalid path",
			})
			return
		}
		c.JSON(201, pkiResponseJson)
	}
}
