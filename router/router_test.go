package router

import (
	"context"
	"testing"

	"github.com/0gener/banking-core-gateway/middleware"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type testJwtMiddleware struct{}

func (j *testJwtMiddleware) EnsureValidToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := &middleware.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				Subject: "userId|1234",
			},
		}
		r := c.Request.Clone(context.WithValue(context.Background(), jwtmiddleware.ContextKey{}, claims))
		c.Request = r
	}
}

func TestNew(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(&testJwtMiddleware{}, nil)
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
