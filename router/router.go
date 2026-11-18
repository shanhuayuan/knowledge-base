package router

import (
	"app/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	object := api.Group("/object")
	{
		object.GET("/list", handler.ListObjectsHandler)
		object.GET("/info", handler.GetObjectInfoHandler)
		object.DELETE("/delete", handler.DeleteObjectHandler)
		object.POST("/upload", handler.UploadObjectHandler)
		object.POST("/bucket/create", handler.CreateBucketHandler)
	}

	api.POST("/embed", handler.EmbedHandler)

	database := api.Group("/database")
	{
		table := database.Group("table")
		table.POST("/create", handler.CreateTableHandler)
		table.DELETE("/:tableName/delete", handler.DropTableHandler)
		table.GET("/list", handler.ListTablesHandler)
	}

	provider := api.Group("/provider")
	{
		provider.GET("/list", handler.ListProvidersHandler)
		provider.POST("/create", handler.CreateProviderHandler)
		provider.DELETE("/:providerId/delete", handler.DeleteProviderHandler)
		provider.GET("/:providerId/info", handler.GetProviderHandler)
	}

	model := provider.Group("/model")
	{
		model.POST("/create", handler.CreateModelHandler)
		model.GET("/list", handler.ListModelsHandler)
		model.DELETE("/:modelId/delete", handler.DeleteModelHandler)
	}

	vectorDb := api.Group("/vectordb")
	{
		vectorDb.POST("/search", handler.SearchHandler)
	}

	collection := vectorDb.Group("/collection")
	{
		collection.POST("/create", handler.CreateCollectionHandler)
		collection.DELETE("/:collectionName/drop", handler.DropCollectionHandler)
		collection.POST("/:collectionName/load", handler.LoadCollectionHandler)
		collection.GET("/:collectionName/info", handler.DescribeCollectionHandler)
		collection.GET("/list", handler.ListCollectionsHandler)
	}

	knowledge := api.Group("/knowledge")
	{
		// 知识库
		knowledge.POST("/create", handler.CreateKnowledgeBaseHandler)
		knowledge.DELETE("/:kbId/delete", handler.DeleteKnowledgeBaseHandler)
		knowledge.POST("/:kbId/query", handler.QueryHandler)
		knowledge.GET("/:kbId/info", handler.GetKnowledgeBaseHandler)
		knowledge.GET("/:kbId/existed", handler.IsKnowledgeBaseExistHandler)
		knowledge.GET("/list", handler.ListKnowledgeBaseHandler)

		// 文档
		documents := knowledge.Group("/:kbId/documents")
		{
			documents.POST("/upload", handler.UploadDocumentHandler)
			documents.DELETE("/:docId/delete", handler.DeleteDocumentHandler)
			documents.GET("/:docId/info", handler.GetDocumentHandler)
			documents.GET("/:docId/chunk", handler.ChunkDocumentHandler)
			documents.POST("/:docId/parse", handler.ParseDocumentHandler)
			documents.GET("/list", handler.ListDocumentsHandler)
		}

		chunks := documents.Group("/chunks")
		{
			chunks.POST("/insert", handler.InsertChunksHandler)
			chunks.DELETE("/:chunkId/delete", handler.DeleteChunkHandler)
		}
	}

	chunker := api.Group("/chunker")
	{
		chunker.POST("/split/default", handler.SplitHandler)
	}

	return router
}
