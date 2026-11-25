package handler

import (
	"app/global"
	"app/handler/api"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shanhuayuan/object"
)

func ListObjectsHandler(ctx *gin.Context) {
	var o object.Object

	if err := ctx.ShouldBindBodyWithJSON(&o); err != nil {
		api.BadRequestError(ctx, "")
		return
	}

	if o.BucketName == "" || o.Prefix == "" {
		api.BadRequestError(ctx, "存储桶或者路径不能为空")
		return
	}

	result := global.OverallRuntime.ObjectStore.ListObjects(ctx, o)
	api.Ok(ctx, result)
}

func DeleteObjectHandler(ctx *gin.Context) {
	var o object.Object

	if err := ctx.ShouldBindBodyWithJSON(&o); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := global.OverallRuntime.ObjectStore.DeleteObject(ctx, o); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func GetObjectInfoHandler(ctx *gin.Context) {
	var o object.Object

	if err := ctx.ShouldBindBodyWithJSON(&o); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	info, err := global.OverallRuntime.ObjectStore.StatObjectInfo(ctx, o)

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

	if !strings.HasSuffix(uploadObjectOption.Prefix, "/") {
		api.BadRequestError(ctx, "请输入正确的对象存储的位置信息！")
		return
	}

	file, header, err := ctx.Request.FormFile("file")

	if err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	defer file.Close()

	uploadObjectOption.File = file
	uploadObjectOption.Prefix = filepath.Join(uploadObjectOption.Prefix, header.Filename)
	uploadObjectOption.ContentType = header.Header.Get("Content-Type")

	if err := global.OverallRuntime.ObjectStore.UploadObject(ctx, uploadObjectOption); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	result := object.Object{
		BucketName: uploadObjectOption.BucketName,
		Prefix:     uploadObjectOption.Prefix,
	}

	api.Ok(ctx, result)
}

func CreateBucketHandler(ctx *gin.Context) {
	var bucket object.Bucket

	if err := ctx.ShouldBindBodyWithJSON(&bucket); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if _, err := global.OverallRuntime.ObjectStore.CreateBucket(ctx, bucket.Name); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}
