package router

import (
	"context"
	"net/http"

	"github.com/0gener/banking-core-accounts/proto"
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

func (c *accountsController) createAccountHandler(ctx *gin.Context) {
	req := createAccountRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	if req.Currency == nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := c.accountsClient.CreateAccount(context.Background(), &proto.CreateAccountRequest{
		UserId:   "1234",
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
