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
)

var mutexPhysicalFile = &sync.Mutex{}

type FileStorageConfig struct {
	Path string `json:"path"`
}

const kvDataJsonFileName string = "kvdata.json"

func keyToLowerCase(key *string) {
	strings.ToLower(*key)
}

func (fs FileStorageConfig) ListKeys(key *string) ([]string, error) {
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	keyToLowerCase(key)
	keys := make([]string, 0)
	rootKeyPath := path.Join((*(KeyValuePersistence.(*FileStorageConfig))).Path, *key)
	files, err := os.ReadDir(rootKeyPath)
	if err != nil {
		return nil, errors.New("error while reading dir: " + rootKeyPath + "\n" + err.Error())
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
	keyPath := path.Join((*(KeyValuePersistence.(*FileStorageConfig))).Path, *key)
	versions := make([]int, 0)
	dirInfo, err := os.Stat(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		return nil, errors.New("error while getting dir info for key: " + *key + "\n" + err.Error())
	}
	if !dirInfo.IsDir() {
		return nil, errors.New("key is not a directory, key: " + *key)
	}
	files, err := os.ReadDir(keyPath)
	if err != nil {
		return nil, errors.New("error while listing versions for key: " + *key + "\n" + err.Error())
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
	var allVersions, err = (*(KeyValuePersistence.(*FileStorageConfig))).ListVersions(key)
	if err != nil {
		return 0, errors.New("error while getting latest version for key: " + *key + "\n" + err.Error())
	}
	if len(allVersions) == 0 {
		return 0, nil
	}
	return allVersions[len(allVersions)-1], nil
}

func (fs FileStorageConfig) GetNextVersion(key *string) (int, error) {
	var allVersions, err = (*(KeyValuePersistence.(*FileStorageConfig))).ListVersions(key)
	if err != nil {
		return 0, errors.New("error while getting next version for key: " + *key + "\n" + err.Error())
	}
	if len(allVersions) == 0 {
		return 1, nil
	}
	return allVersions[len(allVersions)-1] + 1, nil
}

func (fs FileStorageConfig) Get(key *string, version *int) (*KVData, error) {
	keyToLowerCase(key)
	if version == nil {
		nextVersion, err := (*(KeyValuePersistence.(*FileStorageConfig))).GetLatestVersion(key)
		if err != nil {
			return nil, errors.New("error while getting the key: " + *key + "\n" + err.Error())
		}
		if nextVersion == 0 {
			return nil, nil
		} else {
			version = &nextVersion
		}
	}
	log.Println("Getting key: ", *key, ", version:", *version)
	fullPath := path.Join((*(KeyValuePersistence.(*FileStorageConfig))).Path, *key, strconv.Itoa(*version), kvDataJsonFileName)
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	d, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, errors.New("error while reading file: " + fullPath + "\n" + err.Error())
	}
	var kvData KVData
	errUnmarshal := json.Unmarshal(d, &kvData)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return &kvData, nil
}

func (fs FileStorageConfig) Save(key *string, keyValue *KVData) error {
	latestVersion, err := (*(KeyValuePersistence.(*FileStorageConfig))).GetLatestVersion(key)
	if err != nil {
		return errors.New("error while getting latest version while saving key: " + *key + "\n" + err.Error())
	}
	if latestVersion > 0 {
		return errors.New("version already exists, use Update")
	}
	return (*(KeyValuePersistence.(*FileStorageConfig))).saveOrUpdate(key, keyValue, nil)
}

func (fs FileStorageConfig) Update(key *string, keyValue *KVData, version *int) error {
	latestVersion, err := (*(KeyValuePersistence.(*FileStorageConfig))).GetLatestVersion(key)
	if err != nil {
		return errors.New("error while getting latest version while updating key: " + *key + "\n" + err.Error())
	}
	if latestVersion == 0 {
		return errors.New("version does not exist, use Save")
	}
	return (*(KeyValuePersistence.(*FileStorageConfig))).saveOrUpdate(key, keyValue, version)
}

func (fs FileStorageConfig) saveOrUpdate(key *string, keyValue *KVData, version *int) error {
	log.Println("Save or Update called" + *key + " " + keyValue.Value)
	keyToLowerCase(key)
	if version == nil {
		nextVersion, err := (*(KeyValuePersistence.(*FileStorageConfig))).GetNextVersion(key)
		if err != nil {
			return errors.New("error while getting next version while save or update key: " + *key + "\n" + err.Error())
		}
		version = &nextVersion
	}
	kvDataFilePath := path.Join((*(KeyValuePersistence.(*FileStorageConfig))).Path, *key, strconv.Itoa(*version), kvDataJsonFileName)
	keyVersionDirPath := filepath.Dir(kvDataFilePath)
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	errMakeDir := os.MkdirAll(keyVersionDirPath, os.FileMode(0700))
	if errMakeDir != nil {
		return errors.New("error while creating dir: " + keyVersionDirPath + "\n" + errMakeDir.Error())
	}

	d, errMarshal := json.Marshal(keyValue)

	if errMarshal != nil {
		return errMarshal
	}
	errWriteFile := os.WriteFile(kvDataFilePath, d, os.FileMode(0700))
	if errWriteFile != nil {
		return errWriteFile
	}
	return nil
}

func (fs FileStorageConfig) Delete(key *string, version *int) error {
	keyToLowerCase(key)
	if version == nil {
		latestVersion, err := (*(KeyValuePersistence.(*FileStorageConfig))).GetLatestVersion(key)
		if err != nil {
			return errors.New("error while getting latest version while deleting key: " + *key + "\n" + err.Error())
		}
		if latestVersion == 0 {
			return errors.New("no version exists to delete for key: " + *key)
		}
		version = &latestVersion
	}
	p := path.Join((*(KeyValuePersistence.(*FileStorageConfig))).Path, *key, strconv.Itoa(*version))
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	err := os.Remove(p)
	return err
}
