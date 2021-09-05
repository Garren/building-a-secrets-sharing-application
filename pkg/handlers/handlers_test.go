package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVerbSecretHandler(t *testing.T) {
	verbtests := []struct {
		verb string
		code int
		body string
	}{
		{"GET", http.StatusNotFound, `{"data":""}`},
		{"POST", http.StatusBadRequest, ""},
		{"PUT", http.StatusMethodNotAllowed, "method not allowed"},
		{"DELETE", http.StatusMethodNotAllowed, "method not allowed"},
	}

	mux := http.NewServeMux()
	SetupHandlers(mux)

	for _, tt := range verbtests {
		t.Run(tt.verb, func(t *testing.T) {
			writer := httptest.NewRecorder()
			request, _ := http.NewRequest(tt.verb, "/", bytes.NewReader(nil))
			mux.ServeHTTP(writer, request)
			if writer.Code != tt.code {
				t.Errorf("Response code is %v, expecting %v",
					writer.Code, http.StatusBadRequest)
			}
			body, _ := io.ReadAll(writer.Body)
			if strings.TrimRight(string(body), "\n") != tt.body {
				t.Errorf("Response body not ok '%s'", string(body))
			}
		})
	}
}

func TestVerbHealthCheck(t *testing.T) {
	verbtests := []struct {
		verb string
		code int
		body string
	}{
		{"GET", http.StatusOK, `ok`},
		{"POST", http.StatusMethodNotAllowed, "method not allowed"},
		{"PUT", http.StatusMethodNotAllowed, "method not allowed"},
		{"DELETE", http.StatusMethodNotAllowed, "method not allowed"},
	}

	mux := http.NewServeMux()
	SetupHandlers(mux)

	for _, tt := range verbtests {
		t.Run(tt.verb, func(t *testing.T) {
			writer := httptest.NewRecorder()
			request, _ := http.NewRequest(tt.verb, "/healthcheck", bytes.NewReader(nil))
			mux.ServeHTTP(writer, request)
			if writer.Code != tt.code {
				t.Errorf("Response code is %v, expecting %v",
					writer.Code, http.StatusBadRequest)
			}
			body, _ := io.ReadAll(writer.Body)
			if strings.TrimRight(string(body), "\n") != tt.body {
				t.Errorf("Response body not ok '%s'", string(body))
			}
		})
	}
}
