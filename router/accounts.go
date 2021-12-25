package router

import (
	"context"
	"net/http"

	"github.com/0gener/banking-core-accounts/proto"
	"github.com/0gener/banking-core-gateway/middleware"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gin-gonic/gin"
)

type accountsController struct {
	accountsClient proto.AccountsServiceClient
}

func newAccountsController(accountsClient proto.AccountsServiceClient) accountsController {
	return accountsController{
		accountsClient,
	}
}

type createAccountRequest struct {
	Currency *string `json:"currency"`
}

type createAccountResponse struct {
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
}

type getAccountResponse struct {
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
}

func (c *accountsController) createAccountHandler(ctx *gin.Context) {
	var claims = ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*middleware.CustomClaims)

	req := createAccountRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	if req.Currency == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := c.accountsClient.CreateAccount(context.Background(), &proto.CreateAccountRequest{
		UserId:   claims.Subject,
		Currency: *req.Currency,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, createAccountResponse{
		AccountNumber: res.Account.AccountNumber,
		Currency:      res.Account.Currency,
	})
}

func (c *accountsController) getAccountHandler(ctx *gin.Context) {
	var claims = ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*middleware.CustomClaims)

	res, err := c.accountsClient.GetAccount(context.Background(), &proto.GetAccountRequest{
		UserId: claims.Subject,
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if res == nil || res.Account == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, getAccountResponse{
		AccountNumber: res.Account.AccountNumber,
		Currency:      res.Account.Currency,
	})
}
