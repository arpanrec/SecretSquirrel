package storage

import (
	"log"
	"sync"

	"github.com/arpanrec/secureserver/internal/common"
	"github.com/arpanrec/secureserver/internal/encryption"
	"github.com/arpanrec/secureserver/internal/physical"
)

var pStorage physical.Storage

var once = &sync.Once{}

func getStorage() physical.Storage {
	once.Do(func() {
		storageType := common.GetConfig()["storage"].(map[string]interface{})["type"].(string)
		log.Print("Storage type set to ", storageType)
		switch storageType {
		case "file":
			pStorage = physical.FileStorageConfig{}
		default:
			log.Println("Error Invalid storage type ", storageType)
		}
	})
	return pStorage
}

func GetData(l string) (string, error) {
	s := getStorage()
	d, err := s.GetData(l)
	if err != nil {
		return "", err
	}
	e := encryption.DecryptMessage(&d)
	if e != nil {
		log.Println("Error while getting data while decrypting message: ", e)
		return "", e
	}
	return d, nil
}

func PutData(l string, d string) (bool, error) {
	s := getStorage()
	err := encryption.EncryptMessage(&d)
	if err != nil {
		return false, err
	}
	return s.PutData(l, d)
}

func DeleteData(l string) error {
	s := getStorage()
	return s.DeleteData(l)
}
