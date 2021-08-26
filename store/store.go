package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Store struct {
	path  string
	mutex *sync.Mutex
	data  map[string]string
}

func init() {
	path := os.Getenv("DATA_FILE_PATH")
	if path == "" {
		err := errors.New("DATA_FILE_PATH variable does not exist")
		panic(err)
	}

	if _, err := os.Stat(path); err == nil {
		return
	} else if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		panic(err)
	}
}

func NewStore() Store {
	return Store{
		path:  os.Getenv("DATA_FILE_PATH"),
		mutex: &sync.Mutex{},
		data:  make(map[string]string),
	}
}

func (s *Store) Load() (err error) {

	jsonFile, err := os.Open(s.path)
	if err != nil {
		fmt.Println("Error opening JSON file", err)
		return err
	}
	defer jsonFile.Close()

	if jsonData, err := ioutil.ReadAll(jsonFile); err == nil {
		err = json.Unmarshal(jsonData, &s.data)
	}

	return
}

func (s *Store) write() (err error) {

	jsonFile, err := os.Open(s.path)
	if err != nil {
		fmt.Println("Error opening JSON file", err)
		return err
	}
	defer jsonFile.Close()

	result, err := json.Marshal(s.data)
	err = ioutil.WriteFile(s.path, result, 0644)

	return
}

func (s *Store) DeleteAndWrite(key string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
	err = s.write()

	return
}

func (s *Store) StoreAndWrite(key, value string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
	err = s.write()

	return
}

func (s *Store) Remove(key string) (val string, ok bool, err error) {
	if val, ok = s.data[key]; ok {
		err = s.DeleteAndWrite(key)
	}
	return val, ok, err
}
