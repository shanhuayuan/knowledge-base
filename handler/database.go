package handler

import (
	"app/database"
	"app/global"
	"app/handler/api"
	"app/model"

	"github.com/gin-gonic/gin"
)

func ListTablesHandler(ctx *gin.Context) {
	tables, err := global.OverallRuntime.DataStore.ListTables()

	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, tables)
}

func CreateTableHandler(ctx *gin.Context) {
	if err := global.OverallRuntime.DataStore.CreateTable(&model.Model{}); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func DropTableHandler(ctx *gin.Context) {
	var option database.DropTableOption

	if err := ctx.ShouldBindUri(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.DataStore.DropTable(option.TableName); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}
