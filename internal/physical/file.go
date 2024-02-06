package physical

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/arpanrec/secretsquirrel/internal/appconfig"
)

var mutexPhysicalFile = &sync.Mutex{}

var oncePhysicalFile = &sync.Once{}

type FileStorageConfig struct {
	Path string `json:"path"`
}

const kvDataJsonFileName string = "kvdata.json"

var fileStorageConfigVar FileStorageConfig

func getFileStorageConfigVar() FileStorageConfig {
	oncePhysicalFile.Do(func() {
		storagePath := appconfig.GetConfig().Storage.Config["path"].(string)
		if storagePath == "" {
			log.Fatalln("Fatal Storage path not set")
		}
		absolutePath, err := filepath.Abs(storagePath)
		if err != nil {
			log.Fatalln("Fatal Storage path not valid", err)
		}
		fileStorageConfigVar = FileStorageConfig{
			Path: absolutePath,
		}
		log.Printf("File storage path set to %v", fileStorageConfigVar)
	})
	return fileStorageConfigVar
}

func keyToLowerCase(key *string) {
	strings.ToLower(*key)
}

func (fs FileStorageConfig) ListKeys(key *string) ([]string, error) {
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	keyToLowerCase(key)
	keys := make([]string, 0)
	rootKeyPath := path.Join(getFileStorageConfigVar().Path, *key)
	files, err := os.ReadDir(rootKeyPath)
	if err != nil {
		log.Println("Error while reading dir", err)
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			keys = append(keys, file.Name())
		}
	}
	return keys, nil
}

func (fs FileStorageConfig) ListVersions(key *string) ([]int, error) {
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	keyToLowerCase(key)
	keyPath := path.Join(getFileStorageConfigVar().Path, *key)
	versions := make([]int, 0)
	dirInfo, err := os.Stat(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		log.Println("Error while getting dir info for key", *key, err)
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, errors.New("key is not a directory, key: " + *key)
	}
	files, err := os.ReadDir(keyPath)
	if err != nil {
		log.Println("Error while listing versions for key", *key, err)
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			fileNameAsInt, err := strconv.Atoi(file.Name())
			if err != nil {
				log.Println("Error while converting file name to int", file.Name(), err)
				continue
			}
			versions = append(versions, fileNameAsInt)
		}
	}
	log.Println("Versions for key", *key, "are", versions)
	sort.Ints(versions)
	log.Println("Versions for key", *key, "are", versions)
	return versions, nil
}

func (fs FileStorageConfig) GetLatestVersion(key *string) (int, error) {
	var allVersions, err = getFileStorageConfigVar().ListVersions(key)
	if err != nil {
		log.Println("Error while getting latest version for key", *key, err)
		return 0, err
	}
	if len(allVersions) == 0 {
		return 0, nil
	}
	return allVersions[len(allVersions)-1], nil
}

func (fs FileStorageConfig) GetNextVersion(key *string) (int, error) {
	var allVersions, err = getFileStorageConfigVar().ListVersions(key)
	if err != nil {
		return 0, err
	}
	if len(allVersions) == 0 {
		return 1, nil
	}
	return allVersions[len(allVersions)-1] + 1, nil
}

func (fs FileStorageConfig) Get(key *string, version *int) (*KVData, error) {
	keyToLowerCase(key)
	if version == nil {
		nextVersion, err := getFileStorageConfigVar().GetLatestVersion(key)
		if err != nil {
			return nil, err
		}
		if nextVersion == 0 {
			return nil, nil
		} else {
			version = &nextVersion
		}
	}
	log.Println("Getting key: ", *key, ", version:", *version)
	fullPath := path.Join(getFileStorageConfigVar().Path, *key, strconv.Itoa(*version), kvDataJsonFileName)
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	d, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	var kvData KVData
	errUnmarshal := json.Unmarshal(d, &kvData)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return &kvData, nil
}

func (fs FileStorageConfig) Save(key *string, keyValue *KVData) error {
	latestVersion, err := getFileStorageConfigVar().GetLatestVersion(key)
	if err != nil {
		log.Println("Error while getting latest version while saving key", *key, err)
		return err
	}
	if latestVersion > 0 {
		return errors.New("version already exists, use Update")
	}
	return getFileStorageConfigVar().saveOrUpdate(key, keyValue, nil)
}

func (fs FileStorageConfig) Update(key *string, keyValue *KVData, version *int) error {
	latestVersion, err := getFileStorageConfigVar().GetLatestVersion(key)
	if err != nil {
		return err
	}
	if latestVersion == 0 {
		return errors.New("version does not exist, use Save")
	}
	return getFileStorageConfigVar().saveOrUpdate(key, keyValue, version)
}

func (fs FileStorageConfig) saveOrUpdate(key *string, keyValue *KVData, version *int) error {
	log.Println("Save or Update called" + *key + " " + keyValue.Value)
	keyToLowerCase(key)
	if version == nil {
		nextVersion, err := getFileStorageConfigVar().GetNextVersion(key)
		if err != nil {
			log.Println("Error while getting next version while save or update key", *key, err)
			return err
		}
		version = &nextVersion
	}
	kvDataFilePath := path.Join(getFileStorageConfigVar().Path, *key, strconv.Itoa(*version), kvDataJsonFileName)
	keyVersionDirPath := filepath.Dir(kvDataFilePath)
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	errMakeDir := os.MkdirAll(keyVersionDirPath, 0755)
	if errMakeDir != nil {
		log.Println("Error while creating dir", keyVersionDirPath, errMakeDir)
		return errMakeDir
	}

	d, errMarshal := json.Marshal(keyValue)

	if errMarshal != nil {
		return errMarshal
	}
	errWriteFile := os.WriteFile(kvDataFilePath, d, 0644)
	if errWriteFile != nil {
		return errWriteFile
	}
	return nil
}

func (fs FileStorageConfig) Delete(key *string, version *int) error {
	keyToLowerCase(key)
	if version == nil {
		latestVersion, err := getFileStorageConfigVar().GetLatestVersion(key)
		if err != nil {
			return err
		}
		if latestVersion == 0 {
			return errors.New("no version exists to delete for key: " + *key)
		}
		version = &latestVersion
	}
	p := path.Join(getFileStorageConfigVar().Path, *key, strconv.Itoa(*version))
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	err := os.Remove(p)
	return err
}
