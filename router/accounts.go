package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	ClientId *string `json:"client_id"`
}

type createAccountResponse struct {
	Message string `json:"message"`
}

func createAccountHandler(ctx *gin.Context) {
	body := createAccountRequest{}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}

	ctx.JSON(http.StatusCreated, createAccountResponse{
		Message: "account was created",
	})
}
