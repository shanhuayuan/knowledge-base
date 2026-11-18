package vectordb

import (
	"context"
	"errors"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type InsertOption[T any] struct {
	CollectionName string `json:"collectionName"`
	Entries        []T    `json:"entries"`
}

type DeleteOption struct {
	CollectionName string `json:"collectionName"`
	Expr           string `json:"expr"`
}

func Insert[T any](ctx context.Context, s Store, option InsertOption[T]) error {
	if s == nil {
		return errors.New("store is nil")
	}

	if len(option.Entries) == 0 {
		return errors.New("option entries is empty")
	}

	store, ok := s.(*MilvusStore)
	if !ok {
		return errors.New("store is not milvus store")
	}

	rows := make([]any, 0, len(option.Entries))
	for _, entry := range option.Entries {
		rows = append(rows, entry)
	}

	insertOption := milvusclient.NewRowBasedInsertOption(option.CollectionName, rows...)

	if _, err := store.Client.Insert(ctx, insertOption); err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, s Store, option DeleteOption) error {
	deleteOption := milvusclient.NewDeleteOption(option.CollectionName).WithExpr(option.Expr)

	store, ok := s.(*MilvusStore)
	if !ok {
		return errors.New("store is not milvus store")
	}

	if _, err := store.Client.Delete(ctx, deleteOption); err != nil {
		return err
	}

	return nil
}
