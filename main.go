package main

import (
	"net/http"

	"git.sr.ht/~garren/milestone1-code/controller"
	"git.sr.ht/~garren/milestone1-code/handler"
	"git.sr.ht/~garren/milestone1-code/store"
)

func main() {
	server := http.Server{
		Addr: "0.0.0.0:8080",
	}
	s := store.NewStore()
	c := controller.NewController(s)
	http.HandleFunc("/healthcheck", handler.HealthCheckHandler)
	http.HandleFunc("/", c.SecretHandler)
	server.ListenAndServe()
}
