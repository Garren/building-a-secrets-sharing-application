package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Garren/building-a-secrets-sharing-application/store"
	"github.com/Garren/building-a-secrets-sharing-application/types"
)

func TestGetSecretSuccessPersistence(t *testing.T) {
	id := "7a819afa983d454b3a368c1422ba853c"
	expectedSecret := "My super secret1234151"
	{
		data := map[string]string{id: expectedSecret}
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("failed creating test data %e", err)
		}

		err = ioutil.WriteFile("./testdata/data.json", jsonData, 0644)
		if err != nil {
			t.Errorf("failed creating test data %e", err)
		}
	}
	{
		mux := http.NewServeMux()
		SetupHandlers(mux)
		store.Init("./testdata/data.json")
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
		err := os.Remove("./testdata/data.json")
		if err != nil {
			t.Errorf("failed creating test data %e", err)
		}
	}
}
