package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func BadRequestError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, Response{Message: message})
}

func InternalServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, Response{Message: message})
}

func ConflictError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusConflict, Response{Message: message})
}

func Ok(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{Message: "ok", Data: data})
}

func UnAuthorizedError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{Message: "Unauthorized!", Data: nil})
}

func NotFoundError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, Response{Message: message, Data: nil})
}
