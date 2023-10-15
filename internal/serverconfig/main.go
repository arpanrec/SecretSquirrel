package serverconfig

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var masterServerConfig MasterConfig

var mo = &sync.Once{}
var mu = &sync.Mutex{}

type EncryptionConfig struct {
	GPGPrivateKeyPath string `json:"private_key_path"`
	GPGPublicKeyPath  string `json:"public_key_path"`
	GPGPassphrasePath string `json:"private_key_password_path"`
	DeleteKeys        bool   `json:"delete_key_files_after_startup"`
}

type PkiConfig struct {
	CaCertFile               string `json:"openssl_root_ca_cert_path"`
	CaPrivateKeyFile         string `json:"openssl_root_ca_key_path"`
	CaPrivateKeyPasswordFile string `json:"openssl_root_ca_key_password_path"`
	CaPrivateKeyNopassFile   string `json:"openssl_root_ca_key_nopasswd_path"`
	DeleteKeys               bool   `json:"delete_key_files_after_startup"`
}

type StorageConfig struct {
	StorageType string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
}

type UsersConfig struct {
	UserNames []string `json:"user_names"`
}

type UserConfig struct {
}

type MasterConfig struct {
	Encryption EncryptionConfig `json:"encryption"`
	PkiConfig  PkiConfig        `json:"pki"`
	Storage    StorageConfig    `json:"storage"`
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
		err := json.Unmarshal(configJson, &masterServerConfig)
		if err != nil {
			log.Fatalln("Error unmarshalling config file ", err)
		}
		log.Printf("Config set successfully %v\n", masterServerConfig)
	})
	mu.Unlock()
	return masterServerConfig
}
