package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/arpanrec/secureserver/internal/middleware"
	"github.com/arpanrec/secureserver/internal/routehandlers"
	"github.com/arpanrec/secureserver/internal/serverconfig"
	"github.com/gin-gonic/gin"
)

func ginRunner(serverHosting serverconfig.HostingConfig) {
	if serverHosting.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.JsonLoggerMiddleware())
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalln("Error setting trusted proxies: ", err)
	}
	apiRouter := r.Group("/api")
	log.Println("Starting server on port 8080")
	apiRouterV1 := apiRouter.Group("/v1")
	apiRouterV1.Use(middleware.AuthMiddleWare(), middleware.NameSpaceMiddleWare())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, "LOCK", "UNLOCK"},
		"/tfstate/*any", routehandlers.TfStateHandler())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		"/files/*any", routehandlers.FileServerHandler())
	apiRouterV1.PUT("/pki/*any", routehandlers.PkiHandler())

	if serverHosting.TlsEnable {
		log.Fatal(r.RunTLS("0.0.0.0"+
			":"+strconv.Itoa(serverHosting.Port),
			serverHosting.TlsCertFile,
			serverHosting.TlsKeyFile))
	}
	log.Fatal(r.Run(":" + strconv.Itoa(serverHosting.Port)))
}

func Runner() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverHosting := serverconfig.GetConfig().Hosting
	ginRunner(serverHosting)

}
