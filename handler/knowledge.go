package handler

import (
	"app/document"
	"app/embedding"
	"app/global"
	"app/handler/api"
	"app/knowledge"
	"app/model"
	"app/vectordb"
	"strings"

	"github.com/gin-gonic/gin"
)

// 知识库的增删改查

func CreateKnowledgeBaseHandler(ctx *gin.Context) {
	var kb model.KnowledgeBase
	if err := ctx.ShouldBindBodyWithJSON(&kb); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := knowledge.CreateKnowledgeBase(&kb); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, kb)
}

func DeleteKnowledgeBaseHandler(ctx *gin.Context) {
	var kb model.KnowledgeBase
	if err := ctx.ShouldBindUri(&kb); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := knowledge.DeleteKnowledgeBase(ctx, &kb); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func ListKnowledgeBaseHandler(ctx *gin.Context) {
	var kbs *[]model.KnowledgeBase

	if err := knowledge.ListKnowledgeBase(kbs); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, kbs)
}

func IsKnowledgeBaseExistHandler(ctx *gin.Context) {
	var kb model.KnowledgeBase
	if err := ctx.ShouldBindUri(&kb); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	existed, err := knowledge.IsKnowledgeBaseExist(&kb)

	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, existed)
}

func GetKnowledgeBaseHandler(ctx *gin.Context) {
	var kb model.KnowledgeBase

	if err := ctx.ShouldBindUri(&kb); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := knowledge.GetKnowledgeBase(&kb); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, kb)
}

// 文档的增删改查

func UploadDocumentHandler(ctx *gin.Context) {
	var option document.UploadDocumentOption

	if err := ctx.ShouldBind(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	doc, err := document.UploadDocument(ctx, option)
	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, doc)
}

func DeleteDocumentHandler(ctx *gin.Context) {
	var option model.Document
	if err := ctx.ShouldBindUri(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := document.DeleteDocument(ctx, &option); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func GetDocumentHandler(ctx *gin.Context) {
	var option model.Document

	if err := ctx.ShouldBindUri(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := document.GetDocument(&option); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, option)
}

func ListDocumentsHandler(ctx *gin.Context) {
	var option model.Document

	if err := ctx.ShouldBindUri(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	documents, err := document.ListDocuments(option)

	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, documents)
}

func InsertChunksHandler(ctx *gin.Context) {
	var option knowledge.InsertChunksOption

	if err := ctx.ShouldBindBodyWithJSON(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := knowledge.InsertChunks(ctx, option); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func DeleteChunkHandler(ctx *gin.Context) {
	var option knowledge.DeleteChunksOption
	if err := ctx.ShouldBindBodyWithJSON(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := knowledge.DeleteChunk(ctx, option); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func ChunkDocumentHandler(ctx *gin.Context) {
	var option model.Document

	if err := ctx.ShouldBindUri(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	parts, err := document.ChunkDocument(ctx, &option)
	if err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	api.Ok(ctx, parts)
}

func ParseDocumentHandler(ctx *gin.Context) {
	var option model.Document

	if err := ctx.ShouldBindUri(&option); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	if err := document.ProcessDocument(ctx, &option); err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	api.Ok(ctx, nil)
}

func QueryHandler(ctx *gin.Context) {
	var queryOption knowledge.QueryOption

	if err := ctx.ShouldBind(&queryOption); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	kb := &model.KnowledgeBase{Id: queryOption.KbId}

	if err := global.OverallRuntime.DataStore.TakeEntry(kb); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	m := &model.Model{Id: kb.ModelId}

	if err := global.OverallRuntime.DataStore.TakeEntry(m); err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	embedder, err := embedding.GetEmbedder(m.Id)
	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	e, err := embedder.Embed(ctx, []string{queryOption.Query})
	if err != nil {
		api.InternalServerError(ctx, err.Error())
		return
	}

	searchOption := vectordb.SearchOption{
		CollectionName: "kb_" + strings.ReplaceAll(kb.Id, "-", "_"),
		TopK:           10,
		Embeddings:     e,
	}

	result, err := global.OverallRuntime.VectorStore.Search(ctx, searchOption)
	if err != nil {
		api.BadRequestError(ctx, err.Error())
		return
	}

	api.Ok(ctx, result[0])
}
