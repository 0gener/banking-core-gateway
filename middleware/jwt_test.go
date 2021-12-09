package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
)

type AJwtMiddleware struct {
}

func (j *AJwtMiddleware) EnsureValidToken() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func TestSomething(t *testing.T) {
	jwtMiddleware := AJwtMiddleware{}

	jwtMiddleware.EnsureValidToken()
}
