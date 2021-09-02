package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSecret(t *testing.T) {
	mux := http.NewServeMux()
	SetupHandlers(mux)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusNotFound {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.Bytes()
	if strings.TrimRight(string(body), "\n") != `{"data":""}` {
		t.Errorf("Response body not ok '%s'", string(body))
	}
}

func TestPostSecret(t *testing.T) {
	mux := http.NewServeMux()
	SetupHandlers(mux)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(nil))
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusBadRequest {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestBadVerbSecretHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)
	{
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/", bytes.NewReader(nil))
		mux.ServeHTTP(writer, request)
		if writer.Code != http.StatusBadRequest {
			t.Errorf("Response code is %v", writer.Code)
		}
	}
	{
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("PUT", "/", bytes.NewReader(nil))
		mux.ServeHTTP(writer, request)
		if writer.Code != http.StatusMethodNotAllowed {
			t.Errorf("Response code is %v", writer.Code)
		}
	}
	{
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("DELETE", "/", nil)
		mux.ServeHTTP(writer, request)
		if writer.Code != http.StatusMethodNotAllowed {
			t.Errorf("Response code is %v", writer.Code)
		}
	}
}

func TestGetHealthCheck(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.Bytes()
	if string(body) != "ok" {
		t.Errorf("Response body not ok")
	}
}

func TestPostHealthCheck(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusMethodNotAllowed {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestPutHealthCheck(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusMethodNotAllowed {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestDeleteHealthCheck(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/", secretHandler)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/healthcheck", nil)
	mux.ServeHTTP(writer, request)
	if writer.Code != http.StatusMethodNotAllowed {
		t.Errorf("Response code is %v", writer.Code)
	}
}
