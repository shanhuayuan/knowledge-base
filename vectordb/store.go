package vectordb

import (
	"context"

	"github.com/milvus-io/milvus/client/v2/entity"
)

type StoreConfig struct {
	Address string
	Port    string
}

type Field struct {
	Name         string           `json:"name"`
	DataType     entity.FieldType `json:"dataType"`
	IsPrimaryKey bool             `json:"isPrimaryKey"`
	Dim          int64            `json:"dim"`
	AutoId       bool             `json:"autoId,omitempty"`
}

type Collection struct {
	Name   string  `json:"collectionName" uri:"collectionName"`
	Desc   string  `json:"desc"`
	AutoId bool    `json:"autoId"`
	Loaded bool    `json:"loaded"`
	Fields []Field `json:"fields"`
}

type SearchOption struct {
	CollectionName string      `json:"collectionName"`
	TopK           int         `json:"topK"`
	Embeddings     [][]float32 `json:"embeddings"`
}

type SearchResultChunk struct {
	Id      int64   `json:"id"`
	Content string  `json:"content"`
	Score   float32 `json:"score"`
}

type Store interface {
	CreateCollection(ctx context.Context, collection *Collection) error
	DropCollection(ctx context.Context, collectionName string) error
	DescribeCollection(ctx context.Context, collection *Collection) error
	ListCollections(ctx context.Context) ([]string, error)
	HasCollection(ctx context.Context, collectionName string) (bool, error)
	LoadCollection(ctx context.Context, collection *Collection) error
	IsCollectionLoaded(ctx context.Context, name string) (bool, error)
	Search(ctx context.Context, option SearchOption) ([][]SearchResultChunk, error)
}
