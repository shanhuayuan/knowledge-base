package knowledge

import (
	"app/global"
	"app/model"
	"context"
	"fmt"
	"strings"

	"github.com/shanhuayuan/vectordb"
)

type InsertChunksOption struct {
	KbId   string        `json:"kbId"`
	Chunks []model.Chunk `json:"chunks"`
}

type DeleteChunksOption struct {
	KbId  string `json:"kbId"`
	DocId string `json:"docId"`
}

func InsertChunks(ctx context.Context, option InsertChunksOption) error {
	insertOption := vectordb.InsertOption[model.Chunk]{
		CollectionName: "kb_" + strings.ReplaceAll(option.KbId, "-", "_"),
		Entries:        option.Chunks,
	}

	if err := vectordb.Insert(ctx, global.OverallRuntime.VectorStore, insertOption); err != nil {
		return err
	}

	return nil
}

func DeleteChunk(ctx context.Context, option DeleteChunksOption) error {
	deleteOption := vectordb.DeleteOption{
		CollectionName: fmt.Sprintf("kb_%s", strings.ReplaceAll(option.KbId, "-", "_")),
		Expr:           fmt.Sprintf("docId == %q", strings.TrimSpace(option.DocId)),
	}

	if err := vectordb.Delete(ctx, global.OverallRuntime.VectorStore, deleteOption); err != nil {
		return err
	}

	return nil
}
