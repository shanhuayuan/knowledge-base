package model

type Chunk struct {
	Id        string    `json:"id" uri:"chunkId" milvus:"name:id"`
	DocId     string    `json:"docId" milvus:"name:docId"`
	Embedding []float32 `json:"embedding" milvus:"name:embedding"`
	Content   string    `json:"content" milvus:"name:content"`
}
