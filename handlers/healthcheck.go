package handlers

import (
	"io"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "ok")
}
