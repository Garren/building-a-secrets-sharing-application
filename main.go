package main

import (
	"log"
	"net/http"
	"os"

	"git.sr.ht/~garren/milestone1-code/handlers"
	"git.sr.ht/~garren/milestone1-code/store"
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

	store.Init(dataFilePath)

	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		log.Fatalf("server could not start listening on %s. error %v", listenAddr, err)
	}
}
