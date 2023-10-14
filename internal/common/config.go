package common

import (
	"encoding/json"
	"log"
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
		configJson, er := os.ReadFile("config.json")
		if er != nil {
			log.Panicf("Error reading config file %v", er)
		}
		err := json.Unmarshal(configJson, &config)
		if err != nil {
			log.Panicf("Error unmarshalling config file %v", err)
		}
		log.Printf("Config set successfully %v\n", config)
	})
	mu.Unlock()
	return config
}
