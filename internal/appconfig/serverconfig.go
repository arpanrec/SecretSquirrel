package appconfig

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var masterServerConfig *ApplicationMasterConfig

var mo = &sync.Once{}
var mu = &sync.Mutex{}

type ApplicationEncryptionConfig struct {
	GPGPrivateKey           string `json:"gpg_private_key"`
	GPGPublicKey            string `json:"gpg_public_key"`
	GPGPrivateKeyPassphrase string `json:"gpg_private_key_password"`
}

type ApplicationStorageConfig struct {
	StorageType string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
}

type UserConfig struct {
}

type ApplicationServerConfig struct {
	Domain      string `json:"domain"`
	Port        int    `json:"port"`
	TlsEnable   bool   `json:"tls_enabled"`
	TlsCertFile string `json:"tls_cert_file"`
	TlsKeyFile  string `json:"tls_key_file"`
	DebugMode   bool   `json:"debug_mode"`
}

type ApplicationMasterConfig struct {
	Encryption   ApplicationEncryptionConfig `json:"encryption"`
	Storage      ApplicationStorageConfig    `json:"storage"`
	ServerConfig ApplicationServerConfig     `json:"server"`
}

func GetConfig() *ApplicationMasterConfig {
	mu.Lock()
	mo.Do(func() {
		log.Printf("Setting config from %v", "config.json")
		configFilePath := os.Getenv("SECURE_SERVER_CONFIG_FILE_PATH")
		if configFilePath == "" {
			configFilePath = "config.json"
		}
		configJson, er := os.ReadFile(configFilePath)
		if er != nil {
			log.Fatalln("Error reading config file", er)
		}
		log.Println("Config file read successfully : \n", string(configJson))
		var newApplicationMasterConfig ApplicationMasterConfig
		err := json.Unmarshal(configJson, &newApplicationMasterConfig)
		if err != nil {
			log.Fatalln("Error Unmarshal config file ", err)
		}
		masterServerConfig = &newApplicationMasterConfig
		log.Printf("Config set successfully %v\n", *masterServerConfig)
	})
	mu.Unlock()
	return masterServerConfig
}
