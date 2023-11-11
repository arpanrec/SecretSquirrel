package secureserver

import (
	"github.com/arpanrec/secretsquirrel/internal/appconfig"
	"github.com/arpanrec/secretsquirrel/internal/ginhosting"
	"log"
)

func Runner() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverHosting := appconfig.GetConfig().ServerConfig
	ginhosting.GinRunner(serverHosting)

}
