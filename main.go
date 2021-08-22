package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)
type (
	SecretPost struct {
		PlainText string `json:"plain_text"`
	}

	SecretGetResponse struct {
		Data string `json:"data"`
	}

	SecretPostResponse struct {
		Id string `json:"id"`
	}

	Store struct {
		path string
		mutex *sync.Mutex
		data map[string]string
	}

	Controller struct {
		store Store
	}
)

func init() {
	fmt.Fprintln(os.Stdout, "init")
	path := os.Getenv("DATA_FILE_PATH")
	if path == "" {
		err := errors.New("DATA_FILE_PATH variable does not exist")
		panic(err)
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stdout, "file %s exists", path)
		return
	} else if os.IsNotExist(err) {
		file, err := os.Create(path)
		fmt.Fprintf(os.Stdout, "file %s created", path)
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
		path: os.Getenv("DATA_FILE_PATH"),
		mutex: &sync.Mutex{},
		data: make(map[string]string),
	}
}

func (s *Store) write() (err error) {

	// write the map the the file
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

func (s *Store) deleteAndWrite(key string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
	err = s.write()

	return
}

func (s *Store) storeAndWrite(key, value string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
	err = s.write()

	return
}

func (s *Store) remove(key string) (val string, ok bool, err error) {
	if val, ok = s.data[key]; ok {
		err = s.deleteAndWrite(key)
	}
	return val, ok, err
}

func (c *Controller) handleGetSecret(w http.ResponseWriter, r *http.Request) (err error) {

	response := SecretGetResponse{ Data: "" }

	key := path.Base(r.URL.Path)
	if key == "" {
		err := errors.New("key not supplied")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if val, ok, err := c.store.remove(key); ok && err == nil {
		response.Data = val
	} else {
		err = errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	output, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

	return
}

func (c *Controller) handlePostSecret(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var secretPost SecretPost
	json.Unmarshal(body, &secretPost)

	hasher := md5.New()
	hasher.Write([]byte(secretPost.PlainText))
	hash := hex.EncodeToString(hasher.Sum([]byte(nil)))

	err = c.store.storeAndWrite(hash, secretPost.PlainText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := SecretPostResponse{ Id: hash }
	output, err := json.Marshal(&response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

	return
}

func (c *Controller) secretHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = c.handleGetSecret(w, r)
	case "POST":
		err = c.handlePostSecret(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// controller factory function
func NewController(s Store) Controller {
	return Controller{ store: s }
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func main() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}
	s := NewStore()
	c := NewController(s)
	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.HandleFunc("/", c.secretHandler)
	server.ListenAndServe()
}
