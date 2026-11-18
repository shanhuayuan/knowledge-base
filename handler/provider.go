package handler

import (
	"app/handler/api"
	"app/model"
	"app/provider"

	"github.com/gin-gonic/gin"
)

func CreateProviderHandler(ctx *gin.Context) {
	var p model.Provider
	if err := ctx.ShouldBindBodyWithJSON(&p); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := provider.CreateProvider(&p); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, p)
}

func GetProviderHandler(ctx *gin.Context) {
	var p model.Provider

	if err := ctx.ShouldBindUri(&p); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := provider.GetProvider(&p); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, p)
}

func DeleteProviderHandler(ctx *gin.Context) {
	var p model.Provider
	if err := ctx.ShouldBindBodyWithJSON(&p); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := provider.DeleteProvider(&p); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func ListProvidersHandler(ctx *gin.Context) {
	var providers []model.Provider

	if err := provider.ListProviders(&providers); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, providers)
}

func CreateModelHandler(ctx *gin.Context) {
	var m model.Model

	if err := ctx.ShouldBindBodyWithJSON(&m); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := provider.CreateModel(&m); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, m)
}

func DeleteModelHandler(ctx *gin.Context) {
	var m model.Model
	if err := ctx.ShouldBindUri(&m); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := provider.DeleteModel(&m); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func ListModelsHandler(ctx *gin.Context) {
	var models *[]model.Model

	if err := provider.ListModels(models); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, models)
}
