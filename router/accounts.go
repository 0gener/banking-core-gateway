package router

import (
	"context"
	"log"
	"net/http"

	"github.com/0gener/banking-core-accounts/proto"
	"github.com/gin-gonic/gin"
)

type accountsController struct {
	accountsClient proto.AccountsServiceClient
}

func NewAccountsController(accountsClient proto.AccountsServiceClient) accountsController {
	return accountsController{
		accountsClient,
	}
}

type createAccountRequest struct {
	ClientId *string `json:"client_id"`
}

type createAccountResponse struct {
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
}

func (c *accountsController) createAccountHandler(ctx *gin.Context) {
	body := createAccountRequest{}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}

	res, err := c.accountsClient.CreateAccount(context.Background(), &proto.CreateAccountRequest{
		UserId:   "1234",
		Currency: "EUR",
	})

	if err != nil {
		log.Println(err)
	}

	ctx.JSON(http.StatusCreated, createAccountResponse{
		AccountNumber: res.AccountNumber,
		Currency:      res.Currency,
	})
}
