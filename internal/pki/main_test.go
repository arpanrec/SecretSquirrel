package pki

import (
	"testing"
)

func TestCertCa(t *testing.T) {
	// caKeyPath := common.GetConfig()["pki"].(map[string]interface{})["openssl_root_ca_key_path"].(string)
	// caKeyPasswordPath := common.GetConfig()["pki"].(map[string]interface{})["openssl_root_ca_key_password_path"].(string)
	// caNopasswdKeyPath := common.GetConfig()["pki"].(map[string]interface{})["openssl_root_ca_key_nopasswd_path"].(string)
	// log.Printf("Removing password from CA key: %s", caKeyPath)
	// removePassCmd := exec.Command("openssl",
	// 	"rsa",
	// 	"-in", caKeyPath,
	// 	"-passin", "file:"+caKeyPasswordPath,
	// 	"-passout", "pass:\"\"",
	// 	"-out", caNopasswdKeyPath)
	// if err := removePassCmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	// t.Log("Removed password from CA key")
	// var pkiJsonSettings ConfigJsonPki
	// pkiJsonSettingsString := common.GetConfig()["pki"].(string)
	// err := json.Unmarshal([]byte(pkiJsonSettingsString), &pkiJsonSettings)
	// if err != nil {
	// 	log.Fatal("Error unmarshalling pki settings: ", err)
	// }
	// log.Printf("Removing password from CA key: %s", pkiJsonSettings.CaPrivateKeyFile)
	// removePassCmd := exec.Command("openssl",
	// 	"rsa",
	// 	"-in", pkiJsonSettings.CaPrivateKeyFile,
	// 	"-passin", "file:"+pkiJsonSettings.CaPrivateKeyPasswordFile,
	// 	"-passout", "pass:\"\"",
	// 	"-out", pkiJsonSettings.CaPrivateKeyNopassFile)
	// if err := removePassCmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	// t.Log("Removed password from CA key")
}
