package document

import (
	"app/chunk"
	"app/embedding"
	"app/global"
	"app/knowledge"
	"app/model"
	"app/object"
	"app/parser"
	"app/provider"
	"app/vectordb"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
)

type UploadDocumentOption struct {
	KbId      string                `form:"kbId"`
	Name      string                `form:"name"`
	ObjectKey string                `form:"objectKey"`
	File      *multipart.FileHeader `form:"file"`
	FileType  string                `form:"fileType"`
	FileSize  int64                 `form:"fileSize"`
}

func ProcessDocument(ctx context.Context, document *model.Document) error {
	if err := global.OverallRuntime.DataStore.TakeEntry(document); err != nil {
		return err
	}

	option := object.Object{
		BucketName: "knowledge",
		ObjectKey:  document.ObjectKey,
	}

	file, err := global.OverallRuntime.ObjectStore.GetObject(ctx, option)

	if err != nil {
		return err
	}

	defer file.Close()

	txtParser := parser.NewTxtParser()

	result, err := txtParser.Invoke(file)
	if err != nil {
		return err
	}

	if result == nil {
		return errors.New("result is nil")
	}

	chunks, err := global.OverallRuntime.Chunker.Split(result.Content.Text)

	if err != nil {
		return err
	}

	kb := model.KnowledgeBase{Id: document.KbId}
	if err := global.OverallRuntime.DataStore.TakeEntry(&kb); err != nil {
		return err
	}

	collectionName := "kb_" + strings.ReplaceAll(kb.Id, "-", "_")

	existed, err := global.OverallRuntime.VectorStore.HasCollection(ctx, collectionName)

	if err != nil {
		return err
	}

	if !existed {
		m := &model.Model{Id: kb.ModelId}
		if err := provider.GetModel(m); err != nil {
			return err
		}

		createCollectionOption := &vectordb.Collection{
			Name:   collectionName,
			Desc:   kb.Desc,
			AutoId: true,
			Fields: []vectordb.Field{
				{Name: "id", DataType: 5, IsPrimaryKey: true},
				{Name: "docId", DataType: 21, IsPrimaryKey: false},
				{Name: "embedding", DataType: 101, Dim: m.Dim},
				{Name: "content", DataType: 21},
			}}

		if err := global.OverallRuntime.VectorStore.CreateCollection(ctx, createCollectionOption); err != nil {
			return err
		}

		collection := &vectordb.Collection{Name: collectionName}
		if err := global.OverallRuntime.VectorStore.LoadCollection(ctx, collection); err != nil {
			return err
		}
	}

	client, err := embedding.GetEmbedder(kb.ModelId)
	if err != nil {
		return err
	}

	texts := make([]string, 0, len(chunks))

	for _, c := range chunks {
		texts = append(texts, c.Text)
	}

	embeddings, err := client.Embed(ctx, texts)
	if err != nil {
		return err
	}

	bizChunks := make([]model.Chunk, 0, len(embeddings))

	for idx := range len(embeddings) {
		c := model.Chunk{
			DocId:     document.Id,
			Embedding: embeddings[idx],
			Content:   texts[idx],
		}

		bizChunks = append(bizChunks, c)
	}

	insertChunksOption := knowledge.InsertChunksOption{
		KbId:   document.KbId,
		Chunks: bizChunks,
	}

	return knowledge.InsertChunks(ctx, insertChunksOption)
}

func UploadDocument(ctx context.Context, option UploadDocumentOption) (*model.Document, error) {
	// 1. 检查知识库是否存在
	kb := model.KnowledgeBase{Id: option.KbId}

	exist, err := knowledge.IsKnowledgeBaseExist(&kb)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("knowledge base is not exist")
	}

	//上传文件到对象存储

	// 检查桶的存在性
	existed, err := global.OverallRuntime.ObjectStore.ExistBucket(ctx, "knowledge")

	if err != nil {
		return nil, err
	}

	if !existed {
		success, err := global.OverallRuntime.ObjectStore.CreateBucket(ctx, "knowledge")

		if err != nil {
			return nil, err
		}

		if !success {
			return nil, errors.New("failed to create knowledge bucket")
		}
	}

	if option.File == nil {
		return nil, errors.New("file is required!")
	}

	file, err := option.File.Open()
	if err != nil {
		return nil, err
	}

	uploadObjectOption := object.UploadObjectOption{
		BucketName:  "knowledge",
		Prefix:      option.ObjectKey,
		File:        file,
		Size:        option.FileSize,
		ContentType: option.FileType,
	}

	if err := global.OverallRuntime.ObjectStore.UploadObject(ctx, uploadObjectOption); err != nil {
		return nil, err
	}

	// 数据库添加记录
	document := &model.Document{
		KbId:       option.KbId,
		Name:       option.Name,
		FileType:   option.FileType,
		FileSize:   option.FileSize,
		ObjectKey:  option.ObjectKey,
		Status:     model.DocumentNotStarted,
		ChunkCount: 0,
	}

	if err := global.OverallRuntime.DataStore.CreateEntry(document); err != nil {
		return nil, err
	}

	return document, nil
}

func DeleteDocument(ctx context.Context, document *model.Document) error {
	if err := global.OverallRuntime.DataStore.TakeEntry(document); err != nil {
		return err
	}

	option := object.Object{
		BucketName: "knowledge",
		ObjectKey:  document.ObjectKey,
	}

	// 从对象存储中删除文件
	if err := global.OverallRuntime.ObjectStore.DeleteObject(ctx, option); err != nil {
		return err
	}

	// 删除数据库内的实体信息
	if err := global.OverallRuntime.DataStore.DeleteEntry(document); err != nil {
		return err
	}

	// 删除向量数据库内的片段信息
	deleteOption := vectordb.DeleteOption{
		CollectionName: fmt.Sprintf("kb_%s", strings.ReplaceAll(document.KbId, "-", "_")),
		Expr:           fmt.Sprintf("docId == %q", strings.TrimSpace(document.Id)),
	}
	if err := vectordb.Delete(ctx, global.OverallRuntime.VectorStore, deleteOption); err != nil {
		return err
	}

	return nil
}

func GetDocument(document *model.Document) error {
	return global.OverallRuntime.DataStore.TakeEntry(document)
}

func ListDocuments(document model.Document) (*[]model.Document, error) {
	var documents *[]model.Document

	if document.KbId == "" {
		return nil, errors.New("kbId is invalid")
	}

	if err := global.OverallRuntime.DataStore.FindEntry(document, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

func ChunkDocument(ctx context.Context, option *model.Document) ([]chunk.ContentPart, error) {
	var document model.Document

	if option.Id == "" {
		return nil, errors.New("id is invalid")
	}

	if err := global.OverallRuntime.DataStore.TakeEntry(document); err != nil {
		return nil, errors.New("")
	}

	getObjectOption := object.Object{
		BucketName: "knowledge",
		ObjectKey:  document.ObjectKey,
	}

	file, err := global.OverallRuntime.ObjectStore.GetObject(ctx, getObjectOption)
	if err != nil {
		return nil, errors.New("")
	}

	if document.FileType == "text/plain" {
		p := parser.NewTxtParser()
		result, err := p.Invoke(file)
		if err != nil {
			return nil, errors.New("")
		}

		if result == nil {
			return nil, errors.New("result is nil")
		}

		content := result.Content

		chunker := chunk.Chunker{}
		parts, err := chunker.Split(content.Text)
		if err != nil {
			return nil, errors.New("")
		}

		return parts, nil
	}

	return nil, errors.New("")
}
