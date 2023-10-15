package pki

import (
	"github.com/arpanrec/secureserver/internal/serverconfig"
	"log"
	"os/exec"
	"testing"
)

func TestCertCa(t *testing.T) {
	pkiJsonSettingsString := serverconfig.GetConfig().PkiConfig
	log.Printf("Removing password from CA key: %s", pkiJsonSettingsString.CaPrivateKeyFile)
	log.Println("openssl",
		"rsa",
		"-in", pkiJsonSettingsString.CaPrivateKeyFile,
		"-passin", "file:"+pkiJsonSettingsString.CaPrivateKeyPasswordFile,
		"-passout", "pass:\"\"",
		"-out", pkiJsonSettingsString.CaPrivateKeyNopassFile)
	removePassCmd := exec.Command("openssl",
		"rsa",
		"-in", pkiJsonSettingsString.CaPrivateKeyFile,
		"-passin", "file:"+pkiJsonSettingsString.CaPrivateKeyPasswordFile,
		"-passout", "pass:\"\"",
		"-out", pkiJsonSettingsString.CaPrivateKeyNopassFile)
	if err := removePassCmd.Run(); err != nil {
		log.Fatal(err)
	}
	t.Log("Removed password from CA key")
}
