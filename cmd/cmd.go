package cmd

import (
	"log"
	"net/http"

	"github.com/arpanrec/secureserver/internal/middleware"
	"github.com/arpanrec/secureserver/internal/routehandlers"
	"github.com/gin-gonic/gin"
)

func Runner() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.JsonLoggerMiddleware())
	apiRouter := r.Group("/api")
	log.Println("Starting server on port 8080")
	apiRouterV1 := apiRouter.Group("/v1")
	apiRouterV1.Use(middleware.AuthMiddleWare(), middleware.NameSpaceMiddleWare())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, "LOCK", "UNLOCK"},
		"/tfstate/*any", routehandlers.TfStateHandler())
	apiRouterV1.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		"/files/*any", routehandlers.FileServerHandler())
	apiRouterV1.PUT("/pki/*any", routehandlers.PkiHandler())
	log.Fatal(r.Run(":8080"))
}
