package embedding

import (
	"app/util"
	"context"
	"errors"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

type OpenAIEmbedding struct {
	Config openai.EmbeddingConfig
}

func NewOpenAIEmbedding(config openai.EmbeddingConfig) *OpenAIEmbedding {
	return &OpenAIEmbedding{Config: config}
}

func (o *OpenAIEmbedding) Embed(ctx context.Context, text []string) ([][]float32, error) {
	embedder, err := openai.NewEmbedder(ctx, &o.Config)

	if err != nil {
		return nil, err
	}

	if embedder == nil {
		return nil, errors.New("embedder is nil")
	}

	embeddings, err := embedder.EmbedStrings(ctx, text)
	if err != nil {
		return nil, err
	}

	return util.Float64MatrixToFloat32(embeddings), nil
}
