package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Garren/building-a-secrets-sharing-application/store"
	"github.com/Garren/building-a-secrets-sharing-application/types"
)

func TestGetSecretSuccessPersistence(t *testing.T) {
	id := "7a819afa983d454b3a368c1422ba853c"
	expectedSecret := "My super secret1234151"
	mux := http.NewServeMux()
	SetupHandlers(mux)
	store.Init("./testdata/data.json", "one", "two")
	testData := []byte(`{"7a819afa983d454b3a368c1422ba853c":"DzwOZje1bn5O8AG8YPGIPJHTMHR/nfum1HKWtOtobk+i2FLHVDmVTzw8/7AdVydMw5w="}`)
	err := ioutil.WriteFile("./testdata/data.json", testData, 0644)
	if err != nil {
		t.Fatalf("failed creating test data: %e", err)
	}

	{
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/"+id, nil)
		mux.ServeHTTP(writer, request)
		if writer.Code != http.StatusOK {
			t.Errorf("Response code is %v", writer.Code)
		}
		body := writer.Body.Bytes()
		response := types.GetSecretResponse{}
		json.Unmarshal(body, &response)
		if response.Data != expectedSecret {
			t.Errorf("wrong response, expecting %s, got %s",
				expectedSecret, response.Data)
		}
	}

	{
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/"+id, nil)
		mux.ServeHTTP(writer, request)
		if writer.Code != http.StatusNotFound {
			t.Errorf("Response code is %v", writer.Code)
		}
		body := writer.Body.Bytes()
		if strings.TrimRight(string(body), "\n") != `{"data":""}` {
			t.Errorf("wrong response, expecting %s, got %s",
				expectedSecret, `{"data":""}`)
		}
	}

	{
		err := os.Remove("./testdata/data.json")
		if err != nil {
			t.Errorf("failed creating test data %e", err)
		}
	}
}
