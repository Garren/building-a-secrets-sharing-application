package handlers

import (
	"io"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain")
		io.WriteString(w, "ok")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
