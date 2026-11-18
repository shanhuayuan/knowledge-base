package handler

import (
	"app/global"
	"app/handler/api"
	"app/vectordb"

	"github.com/gin-gonic/gin"
)

func ListCollectionsHandler(ctx *gin.Context) {
	collections, err := global.OverallRuntime.VectorStore.ListCollections(ctx)

	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, collections)
}

func LoadCollectionHandler(ctx *gin.Context) {
	var collection vectordb.Collection

	if err := ctx.ShouldBindUri(&collection); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.VectorStore.LoadCollection(ctx, &collection); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func DropCollectionHandler(ctx *gin.Context) {
	var collection vectordb.Collection

	if err := ctx.ShouldBindUri(&collection); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.VectorStore.DropCollection(ctx, collection.Name); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func DescribeCollectionHandler(ctx *gin.Context) {
	var collection vectordb.Collection
	if err := ctx.ShouldBindUri(&collection); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.VectorStore.DescribeCollection(ctx, &collection); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	api.Ok(ctx, collection)
}

func CreateCollectionHandler(ctx *gin.Context) {
	var collection vectordb.Collection

	if err := ctx.ShouldBindBodyWithJSON(&collection); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.VectorStore.CreateCollection(ctx, &collection); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func SearchHandler(ctx *gin.Context) {
	var searchOption vectordb.SearchOption

	if err := ctx.ShouldBindBodyWithJSON(&searchOption); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	result, err := global.OverallRuntime.VectorStore.Search(ctx, searchOption)

	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, result)
}
