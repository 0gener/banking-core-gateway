package router

import (
	"github.com/0gener/banking-core/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func New(jwtMiddleware middleware.JwtMiddleware) *gin.Engine {
	r := gin.Default()

	r.GET("/status" /* , jwtMiddleware.EnsureValidToken() */, getStatusHandler)
	r.POST("/accounts" /* , jwtMiddleware.EnsureValidToken() */, createAccountHandler)

	return r
}
