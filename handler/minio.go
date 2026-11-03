package handler

import (
	"app/global"
	"app/handler/api"
	"app/object"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ListObjectsHandler(ctx *gin.Context) {
	var listObjectsOption object.ListObjectsOption

	if err := ctx.ShouldBind(&listObjectsOption); err != nil {
		api.BadRequestError(ctx, "")
		return
	}

	if listObjectsOption.BucketName == "" || listObjectsOption.Prefix == "" {
		api.BadRequestError(ctx, "存储桶或者路径不能为空")
		return
	}

	result := global.OverallRuntime.Store.ListObjects(ctx, listObjectsOption)
	api.Ok(ctx, result)
}

func DeleteObjectHandler(ctx *gin.Context) {
	var deleteObjectOption object.DeleteObjectOption

	if err := ctx.ShouldBind(&deleteObjectOption); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.Store.DeleteObject(ctx, deleteObjectOption); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func GetObjectInfoHandler(ctx *gin.Context) {
	var getObjectInfoOption object.GetObjectOption

	if err := ctx.ShouldBind(&getObjectInfoOption); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	info, err := global.OverallRuntime.Store.GetObjectInfo(ctx, getObjectInfoOption)

	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	if info == nil {
		api.InternalServerError(ctx, "can not get object info!")
	}

	api.Ok(ctx, info)
}

func UploadObjectHandler(ctx *gin.Context) {
	var uploadObjectOption object.UploadObjectOption

	if err := ctx.ShouldBind(&uploadObjectOption); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	file, header, err := ctx.Request.FormFile("file")

	if err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	defer file.Close()

	uploadObjectOption.File = file
	uploadObjectOption.Prefox = filepath.Join(uploadObjectOption.Prefox, header.Filename)
	uploadObjectOption.ContentType = header.Header.Get("Content-Type")

	if err := global.OverallRuntime.Store.UploadObject(ctx, uploadObjectOption); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}
