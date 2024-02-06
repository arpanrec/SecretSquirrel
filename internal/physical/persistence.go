package physical

import (
	"log"
	"sync"

	"github.com/arpanrec/secretsquirrel/internal/appconfig"
)

var KeyValuePersistence *KeyValueStorage

var once = &sync.Once{}

const InternalStoragePath string = "internal"

func getStorage() *KeyValueStorage {
	once.Do(func() {
		storageConfig := appconfig.GetConfig().Storage
		storageType := storageConfig.StorageType
		log.Print("Storage type set to ", storageType)
		switch storageType {
		case "file":
			var filePersistence KeyValueStorage = FileStorageConfig{
				Path: storageConfig.Config["path"].(string),
			}
			KeyValuePersistence = &filePersistence
		default:
			log.Fatalln("Error Invalid storage type ", storageType)
		}
	})
	return KeyValuePersistence
}

func Get(key *string, version *int) (*KVData, error) {
	log.Println("Get called " + *key)
	s := getStorage()
	d, err := (*s).Get(key, version)
	if err != nil {
		log.Println("Error while getting data: ", err)
		return nil, err
	}
	// err = encryption.DecryptMessage(&d.Value)
	// if err != nil {
	// 	return nil, err
	// }
	return d, nil
}

func Save(key *string, keyValue *KVData) error {
	log.Println("Save called" + *key + " " + keyValue.Value)
	s := getStorage()
	// err := encryption.EncryptMessage(&keyValue.Value)
	// if err != nil {
	// 	return err
	// }
	return (*s).Save(key, keyValue)
}

func Update(key *string, keyValue *KVData, version *int) error {
	s := getStorage()
	// err := encryption.EncryptMessage(&keyValue.Value)
	// if err != nil {
	// 	return err
	// }
	return (*s).Update(key, keyValue, version)
}

func Delete(key *string, version *int) error {
	s := getStorage()
	return (*s).Delete(key, version)
}
