package embedding

import (
	"app/model"
	"app/provider"
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

type EmbedOption struct {
	ModelId string   `json:"modelId"`
	Text    []string `json:"text"`
}

type Store interface {
	Embed(ctx context.Context, text []string) ([][]float32, error)
}

func GetEmbedder(modelId string) (Store, error) {
	m := &model.Model{Id: modelId}

	// 1. 获取向量模型
	if err := provider.GetModel(m); err != nil {
		return nil, err
	}

	// 2. 获取模型厂商
	p := &model.Provider{Id: m.ProviderId}
	if err := provider.GetProvider(p); err != nil {
		return nil, err
	}

	config := openai.EmbeddingConfig{
		APIKey:  p.ApiKey,
		BaseURL: p.BaseUrl,
		Model:   m.Name,
	}
	embedder := NewOpenAIEmbedding(config)
	return embedder, nil
}
