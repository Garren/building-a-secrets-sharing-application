package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/Garren/building-a-secrets-sharing-application/store"
	"github.com/Garren/building-a-secrets-sharing-application/types"
)

func getSecret(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	id = strings.TrimPrefix(id, "/")
	if len(id) == 0 {
		http.Error(w, `{"data":""}`, http.StatusNotFound)
		return
	}
	resp := types.GetSecretResponse{}
	v, err := store.FileStoreConfig.Fs.Read(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(v) == 0 {
		http.Error(w, `{"data":""}`, http.StatusNotFound)
		return
	}

	resp.Data = v
	jd, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(resp.Data) == 0 {
		w.WriteHeader(404)
	}

	w.Write(jd)
}

func getHash(plainText string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(plainText)))
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)

	p := types.CreateSecretPayload{}
	err = json.Unmarshal(bytes, &p)
	if err != nil || len(p.PlainText) == 0 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	digest := getHash(p.PlainText)
	response := types.CreateSecretResponse{Id: digest}

	s := types.SecretData{Id: digest, Secret: p.PlainText}
	err = store.FileStoreConfig.Fs.Write(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getSecret(w, r)
	case "POST":
		createSecret(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
