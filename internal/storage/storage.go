package storage

import (
	"github.com/arpanrec/secureserver/internal/common"
	"github.com/arpanrec/secureserver/internal/physical"
	"log"
	"sync"
)

var pStorage physical.Storage

var once = &sync.Once{}

func getStorage() physical.Storage {
	once.Do(func() {
		storageType := common.GetConfig()["storage"].(map[string]interface{})["type"].(string)
		log.Print("Storage type set to ", storageType)
		switch storageType {
		case "file":
			pStorage = physical.FileStorage{}
		default:
			log.Panicf("Invalid storage type %v\n", storageType)
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
	decryptMessage(&d)
	return d, nil
}

func PutData(l string, d string) (bool, error) {
	s := getStorage()
	encryptMessage(&d)
	return s.PutData(l, d)
}

func DeleteData(l string) error {
	s := getStorage()
	return s.DeleteData(l)
}
