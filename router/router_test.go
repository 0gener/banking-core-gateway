package router

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNew(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(nil, nil)
	ri := r.Routes()

	if len(ri) != 3 {
		t.Errorf("expected 3 routes configured, got %d", len(ri))
	}

	expectedRoutes := []struct {
		method string
		path   string
	}{
		{
			method: "GET",
			path:   "/status",
		},
		{
			method: "POST",
			path:   "/accounts",
		},
		{
			method: "GET",
			path:   "/accounts",
		},
	}

	for _, er := range expectedRoutes {
		found := false
		for _, info := range ri {
			if info.Method == er.method && info.Path == er.path {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected to have configured %s %s", er.method, er.path)
		}
	}
}
