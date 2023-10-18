package serverconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"log"
	"os"
	"sync"
)

var masterServerConfig MasterConfig

var mo = &sync.Once{}
var mu = &sync.Mutex{}

type EncryptionConfig struct {
	GPGPrivateKeyFile       string `json:"gpg_private_key_file"`
	GPGPublicKeyFile        string `json:"gpg_public_key_file"`
	GPGPassphraseFile       string `json:"gpg_private_key_password_file"`
	GPGPrivateKey           string `json:"gpg_private_key"`
	GPGPublicKey            string `json:"gpg_public_key"`
	GPGPrivateKeyPassphrase []byte `json:"gpg_private_key_password"`
	GPGDeleteKeys           bool   `json:"pgp_delete_key_files_after_startup"`
}

type PkiConfig struct {
	CaCertFile                 string `json:"openssl_root_ca_cert_file"`
	CaPrivateKeyFile           string `json:"openssl_root_ca_key_file"`
	CaPrivateKeyPasswordFile   string `json:"openssl_root_ca_key_password_file"`
	CaPrivateKeyNoPasswordFile string `json:"openssl_root_ca_no_password_key_file"`
	CaDeleteKeys               bool   `json:"openssl_delete_key_files_after_startup"`
	CaCert                     *x509.Certificate
	CaPrivateNoPasswordKey     *rsa.PrivateKey
}

type StorageConfig struct {
	StorageType string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
}

type UserConfig struct {
}

type MasterConfig struct {
	Encryption EncryptionConfig      `json:"encryption"`
	PkiConfig  PkiConfig             `json:"pki"`
	Storage    StorageConfig         `json:"storage"`
	UserDb     map[string]UserConfig `json:"users"`
	Hosting    HostingConfig         `json:"server"`
}

type HostingConfig struct {
	Domain      string `json:"domain"`
	Port        int    `json:"port"`
	TlsEnable   bool   `json:"tls_enabled"`
	TlsCertFile string `json:"tls_cert_file"`
	TlsKeyFile  string `json:"tls_key_file"`
}

func GetConfig() MasterConfig {
	mu.Lock()
	mo.Do(func() {
		log.Printf("Setting config from %v", "config.json")
		configFilePath := os.Getenv("SECURE_SERVER_CONFIG_FILE_PATH")
		if configFilePath == "" {
			configFilePath = "/home/clouduser/workspace/secureserver/config.json"
		}
		configJson, er := os.ReadFile(configFilePath)
		if er != nil {
			log.Fatalln("Error reading config file", er)
		}
		log.Println("Config file read successfully : \n", string(configJson))
		err := json.Unmarshal(configJson, &masterServerConfig)
		if err != nil {
			log.Fatalln("Error Unmarshal config file ", err)
		}
		log.Printf("Config set successfully %v\n", masterServerConfig)
	})
	mu.Unlock()
	return masterServerConfig
}
