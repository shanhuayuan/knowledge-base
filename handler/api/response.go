package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func BadRequestError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, ApiResponse{Message: message})
}

func InternalServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, ApiResponse{Message: message})
}

func ConflictError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusConflict, ApiResponse{Message: message})
}

func Ok(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, ApiResponse{Message: "ok", Data: data})
}

func UnAuthorizedError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, ApiResponse{Message: "Unauthorized!", Data: nil})
}

func NotFoundError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, ApiResponse{Message: message, Data: nil})
}
