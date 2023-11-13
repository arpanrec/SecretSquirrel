package secureserver

import (
	"log"

	"github.com/arpanrec/secretsquirrel/internal/appconfig"
	"github.com/arpanrec/secretsquirrel/internal/ginhosting"
)

func Runner() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverHosting := appconfig.GetConfig().ServerConfig
	ginhosting.GinRunner(serverHosting)

}
