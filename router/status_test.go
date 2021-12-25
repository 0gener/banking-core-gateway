package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	r := New(&testJwtMiddleware{}, nil)

	request := httptest.NewRequest("GET", "/status", nil)

	r.ServeHTTP(rec, request)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	bytes, _ := ioutil.ReadAll(rec.Body)
	response := &getStatusResponse{}
	json.Unmarshal(bytes, response)

	if response.Status != "OK" {
		t.Errorf("expected status message %s, got %s", "OK", response.Status)
	}
}
