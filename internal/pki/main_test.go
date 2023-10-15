package pki

import (
	"log"
	"os/exec"
	"testing"

	"github.com/arpanrec/secureserver/internal/serverconfig"
)

func TestCertCa(t *testing.T) {
	pkiJsonSettingsString := serverconfig.GetConfig().PkiConfig
	log.Printf("Removing password from CA key: %s", pkiJsonSettingsString.CaPrivateKeyFile)
	removePassCmd := exec.Command("openssl",
		"rsa",
		"-in", pkiJsonSettingsString.CaPrivateKeyFile,
		"-passin", "file:"+pkiJsonSettingsString.CaPrivateKeyPasswordFile,
		"-passout", "pass:\"\"",
		"-out", pkiJsonSettingsString.CaPrivateKeyNoPasswordFile)
	if err := removePassCmd.Run(); err != nil {
		log.Fatal(err)
	}
	t.Log("Removed password from CA key")
}
