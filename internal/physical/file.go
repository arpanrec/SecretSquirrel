package physical

import (
	"github.com/arpanrec/secureserver/internal/common"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type FileStorage struct{}

var fileStoragePath string

var mutexPhysicalFile = &sync.Mutex{}

var oncePhysicalFile = &sync.Once{}

type FileStorageConfig struct {
	Path string `json:"path"`
}

func getPath() string {
	oncePhysicalFile.Do(func() {
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
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	p := path.Join(getPath(), Location)
	dir := filepath.Dir(p)
	errMakeDir := os.MkdirAll(dir, 0755)
	if errMakeDir != nil {
		log.Println("Error creating directory: ", errMakeDir)
		return false, errMakeDir
	}
	errWriteFile := os.WriteFile(p, []byte(Data), 0644)
	if errWriteFile != nil {
		log.Println("Error writing file: ", errWriteFile)
		return false, errWriteFile
	}
	return true, nil
}

func (fs FileStorage) DeleteData(Location string) error {
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	p := path.Join(getPath(), Location)
	err := os.Remove(p)
	return err
}
