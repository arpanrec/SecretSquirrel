package ginhosting

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arpanrec/secureserver/internal/appconfig"
	"github.com/gin-gonic/gin"
)

func GinRunner(serverHosting appconfig.ApplicationServerConfig) {
	if serverHosting.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(jsonLoggerMiddleware())
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalln("Error setting trusted proxies: ", err)
	}
	apiRouter := r.Group("/api")
	log.Println("Starting server on port 8080")
	apiRouterV1 := apiRouter.Group("/v1")
	apiRouterV1.Use(authMiddleWare(), nameSpaceMiddleWare())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, "LOCK", "UNLOCK"},
		"/tfstate/*any", tfStateHandler())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		"/files/*any", fileServerHandler())
	apiRouterV1.PUT("/pki/*any", pkiHandler())

	if serverHosting.TlsEnable {
		log.Println("Starting server on port", serverHosting.Port,
			" TLS Cert File: ", serverHosting.TlsCertFile,
			" TLS Key File: ", serverHosting.TlsKeyFile)
		log.Fatal(r.RunTLS(":"+strconv.Itoa(serverHosting.Port),
			serverHosting.TlsCertFile,
			serverHosting.TlsKeyFile))
	}
	log.Fatal(r.Run(":" + strconv.Itoa(serverHosting.Port)))
}
