package handler

import (
	"app/embedding"
	"app/handler/api"

	"github.com/gin-gonic/gin"
)

func EmbedHandler(ctx *gin.Context) {
	var option embedding.EmbedOption

	if err := ctx.ShouldBindBodyWithJSON(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	embedder, err := embedding.GetEmbedder(option.ModelId)
	if err != nil {
		return
	}

	result, err := embedder.Embed(ctx, option.Text)
	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, result)
}
