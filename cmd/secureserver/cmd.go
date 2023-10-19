package secureserver

import (
	"github.com/arpanrec/secureserver/internal/appconfig"
	"github.com/arpanrec/secureserver/internal/ginhosting"
	"log"
)

func Runner() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverHosting := appconfig.GetConfig().ServerConfig
	ginhosting.GinRunner(serverHosting)

}
