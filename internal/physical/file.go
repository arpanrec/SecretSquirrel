package physical

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

type Storage interface {
	GetData() (string, error)
	PutData() (bool, error)
	DeleteData() error
}

type FileStorage struct {
	Location string
	Data     string
}

func getPath() string {
	var StorageDataDir = "data"
	issDataDirEnv := os.Getenv("ISS_DATA_DIR")
	if issDataDirEnv != "" {
		StorageDataDir = issDataDirEnv
	}
	return StorageDataDir
}

func (fs FileStorage) GetData() (string, error) {
	p := path.Join(getPath(), fs.Location)

	// Read the file and return the contents
	d, err := os.ReadFile(p)
	return string(d), err
}

func (fs FileStorage) PutData() (bool, error) {
	p := path.Join(getPath(), fs.Location)
	dir := filepath.Dir(p)
	errMakeDir := os.MkdirAll(dir, 0755)
	if errMakeDir != nil {
		log.Fatal(errMakeDir)
		return false, errMakeDir
	}
	errWriteFile := os.WriteFile(p, []byte(fs.Data), 0644)
	if errWriteFile != nil {
		log.Fatal(errWriteFile)
		return false, errWriteFile
	}
	return true, nil
}

func (fs FileStorage) DeleteData() error {
	p := path.Join(getPath(), fs.Location)
	err := os.Remove(p)
	return err
}
