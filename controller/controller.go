package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"git.sr.ht/~garren/milestone1-code/store"
	"git.sr.ht/~garren/milestone1-code/types"
)

type Controller struct {
	store store.Store
}

func (c *Controller) handleGetSecret(w http.ResponseWriter, r *http.Request) (err error) {
	response := types.SecretGetResponse{ Data: "" }

	key := path.Base(r.URL.Path)
	if key == "" {
		err := errors.New("key not supplied")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = c.store.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if val, ok, err := c.store.Remove(key); ok && err == nil {
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

	var secretPost types.SecretPost
	json.Unmarshal(body, &secretPost)

	hasher := md5.New()
	hasher.Write([]byte(secretPost.PlainText))
	hash := hex.EncodeToString(hasher.Sum([]byte(nil)))

	err = c.store.StoreAndWrite(hash, secretPost.PlainText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.SecretPostResponse{ Id: hash }
	output, err := json.Marshal(&response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

	return
}

func (c *Controller) SecretHandler(w http.ResponseWriter, r *http.Request) {
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

func NewController(s store.Store) Controller {
	return Controller{ store: s }
}
