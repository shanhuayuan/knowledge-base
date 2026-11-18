package vectordb

import (
	"context"
	"errors"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type MilvusStore struct {
	Client *milvusclient.Client
}

func (m *MilvusStore) CreateCollection(ctx context.Context, collection *Collection) error {
	has, err := m.HasCollection(ctx, collection.Name)
	if err != nil {
		return err
	}

	if has {
		return errors.New("collection does existed")
	}

	schema := entity.NewSchema().WithName(collection.Name).WithDescription(collection.Desc)

	for _, f := range collection.Fields {
		field := entity.NewField().WithName(f.Name).WithDataType(f.DataType)

		if f.Name == "id" {
			field.WithIsAutoID(collection.AutoId)
		}

		if f.IsPrimaryKey {
			field.WithIsPrimaryKey(true)
		}

		if f.DataType == entity.FieldTypeFloatVector {
			if f.Dim == 0 {
				return errors.New("field dim is invalid")
			}

			field.WithDim(f.Dim)
		}

		if f.DataType == entity.FieldTypeVarChar {
			field.WithMaxLength(65535)
		}

		schema.WithField(field)
	}

	createCollectionOption := milvusclient.NewCreateCollectionOption(collection.Name, schema)

	if err := m.Client.CreateCollection(ctx, createCollectionOption); err != nil {
		return err
	}

	idx := index.NewHNSWIndex(entity.COSINE, 16, 200)
	indexOption := milvusclient.NewCreateIndexOption(collection.Name, "embedding", idx)

	if _, err := m.Client.CreateIndex(ctx, indexOption); err != nil {
		return err
	}

	return nil
}

func (m *MilvusStore) DropCollection(ctx context.Context, collectionName string) error {
	has, err := m.HasCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	if !has {
		return nil
	}

	option := milvusclient.NewDropCollectionOption(collectionName)

	return m.Client.DropCollection(ctx, option)
}

func (m *MilvusStore) ListCollections(ctx context.Context) ([]string, error) {
	option := milvusclient.NewListCollectionOption()
	return m.Client.ListCollections(ctx, option)
}

func (m *MilvusStore) DescribeCollection(ctx context.Context, collection *Collection) error {
	has, err := m.HasCollection(ctx, collection.Name)

	if err != nil || !has {
		return errors.New(collection.Name + " not exist")
	}

	option := milvusclient.NewDescribeCollectionOption(collection.Name)
	entryCollection, err := m.Client.DescribeCollection(ctx, option)

	if err != nil {
		return err
	}

	if entryCollection == nil {
		return errors.New("entryCollection is nil")
	}

	fields := make([]Field, 0)

	for _, field := range entryCollection.Schema.Fields {
		dim, _ := field.GetDim()

		temp := Field{
			Name:         field.Name,
			DataType:     field.DataType,
			IsPrimaryKey: field.PrimaryKey,
			AutoId:       field.AutoID,
			Dim:          dim,
		}

		fields = append(fields, temp)
	}

	collection = &Collection{
		Name:   entryCollection.Name,
		Loaded: entryCollection.Loaded,
		Fields: fields,
	}

	return nil
}

func (m *MilvusStore) HasCollection(ctx context.Context, collectionName string) (bool, error) {
	hasOption := milvusclient.NewHasCollectionOption(collectionName)
	return m.Client.HasCollection(ctx, hasOption)
}

func (m *MilvusStore) IsCollectionLoaded(ctx context.Context, collectionName string) (bool, error) {
	option := milvusclient.NewGetLoadStateOption(collectionName)
	loaded, err := m.Client.GetLoadState(ctx, option)
	return loaded.State == entity.LoadStateLoaded, err
}

func (m *MilvusStore) LoadCollection(ctx context.Context, collection *Collection) error {
	loadOption := milvusclient.NewLoadCollectionOption(collection.Name)
	loadTask, err := m.Client.LoadCollection(ctx, loadOption)

	if err != nil {
		return err
	}

	return loadTask.Await(ctx)
}

func (m *MilvusStore) Search(ctx context.Context, option SearchOption) ([][]SearchResultChunk, error) {
	has, err := m.HasCollection(ctx, option.CollectionName)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("collection not exist")
	}

	loaded, err := m.IsCollectionLoaded(ctx, option.CollectionName)

	if err != nil {
		return nil, err
	}

	if !loaded {
		collection := &Collection{Name: option.CollectionName}
		if err := m.LoadCollection(ctx, collection); err != nil {
			return nil, err
		}
	}

	searchOption := milvusclient.
		NewSearchOption(option.CollectionName, option.TopK, m.data2FloatVectors(option.Embeddings)).
		WithOutputFields("id", "content")

	results, err := m.Client.Search(ctx, searchOption)

	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, errors.New("search result is nil")
	}

	allEntries := make([][]SearchResultChunk, 0, len(results))

	for _, result := range results {
		ids := result.GetColumn("id")
		contents := result.GetColumn("content")
		scores := result.Scores

		count := result.ResultCount

		entries := make([]SearchResultChunk, 0, count)

		for i := range count {
			id, _ := ids.GetAsInt64(i)
			content, _ := contents.GetAsString(i)

			entry := SearchResultChunk{
				Id:      id,
				Content: content,
				Score:   scores[i],
			}
			entries = append(entries, entry)
		}

		allEntries = append(allEntries, entries)
	}

	return allEntries, nil
}

func (m *MilvusStore) data2FloatVectors(matrix [][]float32) []entity.Vector {
	vectors := make([]entity.Vector, 0, len(matrix))

	for _, vector := range matrix {
		vectors = append(vectors, entity.FloatVector(vector))
	}

	return vectors
}
