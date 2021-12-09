package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getStatusResponse struct {
	Status string `json:"status"`
}

func getStatusHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, getStatusResponse{
		Status: "OK",
	})
}
