package handler

import (
	"app/global"
	"app/handler/api"

	"github.com/gin-gonic/gin"
)

type ChunkRequest struct {
	Text string `json:"text"`
}

func SplitHandler(ctx *gin.Context) {
	var req ChunkRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	parts, err := global.OverallRuntime.Chunker.Split(req.Text)
	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, parts)
}
