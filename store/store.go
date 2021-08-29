package store

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"git.sr.ht/~garren/milestone1-code/types"
)

type fileStore struct {
	Store map[string]string
	mutex sync.Mutex
}

var FileStoreConfig struct {
	DataFilePath string
	Fs fileStore
}

func Init(dataFilePath string) error {
	_, err := os.Stat(dataFilePath)

	if err != nil {
		_, err := os.Create(dataFilePath)
		if err != nil {
			return err
		}
	}

	FileStoreConfig.Fs = fileStore{
		mutex: sync.Mutex{},
		Store: make(map[string]string),
	}
	FileStoreConfig.DataFilePath = dataFilePath

	return nil
}

func (j *fileStore) ReadFromFile() error {
	f, err := os.Open(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}
	jsonData, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	if len(jsonData) != 0 {
		return json.Unmarshal(jsonData, &j.Store)
	}
	return nil
}

func (j *fileStore) WriteToFile() error {
	var f *os.File
	jsonData, err := json.Marshal(j.Store)
	if err != nil {
		return err
	}
	f, err = os.Open(FileStoreConfig.DataFilePath)
	if err != nil {
		return err
	}
	jsonData, err = io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	if len(jsonData) != 0 {
		return json.Unmarshal(jsonData, &j.Store)
	}
	return nil
}

func (j *fileStore) Write(data types.SecretData) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := j.ReadFromFile()
	if err != nil {
		return err
	}
	j.Store[data.Id] = data.Secret
	return j.WriteToFile()
}

func (j *fileStore) Read(id string) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := j.ReadFromFile()
	if err != nil {
		return "", err
	}
	data := j.Store[id]
	delete(j.Store, id)
	j.WriteToFile()
	return data, nil
}

