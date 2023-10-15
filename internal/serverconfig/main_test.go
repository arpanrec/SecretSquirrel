package serverconfig

import (
	"github.com/arpanrec/secureserver/internal/physical"
	"log"
	"testing"
)

func TestGetConfig(t *testing.T) {
	got := GetConfig()
	log.Print(got.Storage.Config)
	storageJsonConfig := got.Storage.Config
	stgFile := physical.FileStorageConfig{
		Path: storageJsonConfig["path"].(string),
	}
	log.Printf("storageJsonConfig: %v", stgFile)
}
