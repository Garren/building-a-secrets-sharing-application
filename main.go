package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Garren/building-a-secrets-sharing-application/handlers"
	"github.com/Garren/building-a-secrets-sharing-application/store"
)

func main() {
	listenAddr := ":8080"
	if fromEnv := os.Getenv("LISTEN_ADDR"); fromEnv != "" {
		listenAddr = fromEnv
	}

	mux := http.NewServeMux()
	handlers.SetupHandlers(mux)

	dataFilePath := os.Getenv("DATA_FILE_PATH")
	if len(dataFilePath) == 0 {
		log.Fatal("Specify DATA_FILE_PATH")
	}

	password := os.Getenv("PASSWORD")
	if len(password) == 0 {
		log.Fatal("Specify PASSWORD")
	}

	salt := os.Getenv("SALT")
	if len(salt) == 0 {
		log.Fatal("Specify SALT")
	}

	store.Init(dataFilePath, password, salt)

	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("server could not start listening on %s. error %v", listenAddr, err)
	}
}
