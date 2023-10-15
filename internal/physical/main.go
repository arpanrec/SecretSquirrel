package physical

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/arpanrec/secureserver/internal/common"
)

var mutex = &sync.Mutex{}

var once = &sync.Once{}

type Storage interface {
	GetData(l string) (string, error)
	PutData(l string, d string) (bool, error)
	DeleteData(l string) error
}

type FileStorage struct{}

var fileStoragePath string

func getPath() string {
	once.Do(func() {
		fileStoragePath = common.GetConfig()["storage"].(map[string]interface{})["config"].(map[string]interface{})["path"].(string)
		log.Printf("File storage path set to %v", fileStoragePath)
	})
	return fileStoragePath
}

func (fs FileStorage) GetData(Location string) (string, error) {
	p := path.Join(getPath(), Location)
	d, err := os.ReadFile(p)
	return string(d), err
}

func (fs FileStorage) PutData(Location string, Data string) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()
	p := path.Join(getPath(), Location)
	dir := filepath.Dir(p)
	errMakeDir := os.MkdirAll(dir, 0755)
	if errMakeDir != nil {
		log.Fatalln("Error creating directory: ", errMakeDir)
		return false, errMakeDir
	}
	errWriteFile := os.WriteFile(p, []byte(Data), 0644)
	if errWriteFile != nil {
		log.Fatalln("Error writing file: ", errWriteFile)
		return false, errWriteFile
	}
	return true, nil
}

func (fs FileStorage) DeleteData(Location string) error {
	mutex.Lock()
	defer mutex.Unlock()
	p := path.Join(getPath(), Location)
	err := os.Remove(p)
	return err
}
