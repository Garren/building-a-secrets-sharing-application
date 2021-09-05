package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Garren/building-a-secrets-sharing-application/pkg/types"
)

func urlHelper(url string) string {
	baseURL := url
	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = fmt.Sprintf("http://%s", baseURL)
	}
	return baseURL
}

func createSecret(apiURL string, plainText string) (types.CreateSecretResponse, error) {
	result := types.CreateSecretResponse{}
	payload := types.CreateSecretPayload{PlainText: plainText}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}
	baseURL := urlHelper(apiURL)
	resp, err := http.Post(baseURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("unknown error '%d'", resp.StatusCode))
	} else {
		err = json.Unmarshal(body, &result)
	}
	return result, err
}

func getSecret(apiURL string, secretID string) (types.GetSecretResponse, error) {
	result := types.GetSecretResponse{}
	baseURL := urlHelper(apiURL)
	endpoint := baseURL + "/" + secretID
	resp, err := http.Get(endpoint)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &result)
	} else if resp.StatusCode == http.StatusNotFound {
		err = errors.New("id not found")
	} else {
		err = errors.New(fmt.Sprintf("unknown error '%d'", resp.StatusCode))
	}
	return result, err
}
