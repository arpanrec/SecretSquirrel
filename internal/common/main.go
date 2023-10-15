package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var config map[string]interface{}

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

type MasterConfig struct {
	Encryption EncryptionConfig `json:"encryption"`
	PkiConfig  PkiConfig        `json:"pki"`
}

func GetConfig() map[string]interface{} {
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
		err := json.Unmarshal(configJson, &config)
		if err != nil {
			log.Fatalln("Error unmarshalling config file ", err)
		}
		log.Printf("Config set successfully %v\n", config)
	})
	mu.Unlock()
	return config
}

func DeleteFile(l string) (bool, error) {
	err := os.Remove(l)
	if err != nil {
		log.Println("Error deleting file: ", err)
		return false, err
	}
	return true, nil
}

func HttpResponseWriter(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, body)
	if err != nil {
		log.Println("Error writing response: ", err)
	}
}
