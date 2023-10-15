package serverconfig

import (
	"log"
	"testing"
)

func TestGetConfig(t *testing.T) {
	got := GetConfig()
	log.Print(got.Storage.Config)
	users := got.UserDb
	log.Printf("Users: %v", users)
}
