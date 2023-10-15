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
			log.Println("Error reading config file", er)
		}
		err := json.Unmarshal(configJson, &config)
		if err != nil {
			log.Println("Error unmarshalling config file ", err)
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
