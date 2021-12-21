package router

import (
	"github.com/0gener/banking-core-accounts/proto"
	"github.com/0gener/banking-core-gateway/middleware"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func New(jwtMiddleware middleware.JwtMiddleware, accountsClient proto.AccountsServiceClient) *gin.Engine {
	r := gin.Default()
	accountsController := newAccountsController(accountsClient)
	r.GET("/status", jwtMiddleware.EnsureValidToken(), getStatusHandler)

	r.POST("/accounts", jwtMiddleware.EnsureValidToken(), accountsController.createAccountHandler)
	r.GET("/accounts", jwtMiddleware.EnsureValidToken(), accountsController.getAccountHandler)

	return r
}
