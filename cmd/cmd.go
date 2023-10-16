package cmd

import (
	"github.com/arpanrec/secureserver/internal/middleware"
	"github.com/arpanrec/secureserver/internal/routehandlers"
	"github.com/gin-gonic/gin"
	"log"
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
	apiRouterV1.Any("/tfstate/*any", routehandlers.TfStateHandler())
	apiRouterV1.Any("/files/*any", routehandlers.FileServerHandler())
	log.Fatal(r.Run(":8080"))
}
