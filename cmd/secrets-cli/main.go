package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Garren/building-a-secrets-sharing-application/pkg/types"
)

type Config struct {
	AppUrl    *string
	Action    *string
	PlainText *string
	SecretId  *string
}

func urlHelper(url string) string {
	baseURL := url
	if !strings.HasPrefix(baseURL, "http://") {
		baseURL = fmt.Sprintf("http://%s", baseURL)
	}
	return baseURL
}

func validate(c Config) (err error) {
	msg := ""
	if *c.Action == "create" {
		if len(*c.PlainText) == 0 {
			msg = "create action requires a text argument"
		} else if len(*c.SecretId) > 0 {
			msg = "create action does not accept a secret id"
		}
	}
	if *c.Action == "view" {
		if len(*c.PlainText) > 0 {
			msg = "view action does not accept a text argument"
		} else if len(*c.SecretId) == 0 {
			msg = "view action requires a secret id"
		}
	}
	if len(msg) > 0 {
		err = errors.New(msg)
	}
	return err
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

func main() {
	c := Config{}

	c.AppUrl = flag.String("url", ":8080", "API Endpoint")
	c.Action = flag.String("action", "", "action to perform (create or view a secret)")
	c.PlainText = flag.String("text", "", "secret text (create)")
	c.SecretId = flag.String("id", "", "secret id (get)")

	showHelp := flag.Bool("h", false, "show help")

	flag.Parse()

	if *showHelp {
		fmt.Println("Usage: ...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := validate(c); err != nil {
		log.Fatalf("flag error: %e", err)
	}

	switch *c.Action {
	case "view":
		resp, err := getSecret(*c.AppUrl, *c.SecretId)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(resp.Data)
		}
	case "create":
		resp, err := createSecret(*c.AppUrl, *c.PlainText)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(resp.Id)
		}
	default:
		fmt.Printf("unknown action '%s'", *c.Action)
		os.Exit(1)
	}
}
