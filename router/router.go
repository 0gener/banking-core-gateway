package router

import (
	"github.com/0gener/banking-core-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func New(jwtMiddleware middleware.JwtMiddleware, accountsController accountsController) *gin.Engine {
	r := gin.Default()
	r.GET("/status" /* , jwtMiddleware.EnsureValidToken() */, getStatusHandler)
	r.POST("/accounts" /* , jwtMiddleware.EnsureValidToken() */, accountsController.createAccountHandler)

	return r
}
