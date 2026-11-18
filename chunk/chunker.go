package chunk

import (
	"errors"
	"strings"
)

type ContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Option struct {
	ChunkSize    int      `json:"chunkSize"`
	ChunkOverlap int      `json:"chunkOverlap"`
	Separators   []string `json:"separators"`
}

type Chunker struct{}

func NewTextChunker() *Chunker {
	return &Chunker{}
}

func (c *Chunker) Split(text string) ([]ContentPart, error) {

	content := strings.TrimSpace(text)
	if content == "" {
		return nil, errors.New("content is nil")
	}

	sentences := strings.SplitAfter(content, "。")

	chunks := make([]ContentPart, 0, len(sentences))

	for _, sentence := range sentences {
		if sentence == "" {
			continue
		}

		chunks = append(chunks, ContentPart{Type: "text", Text: sentence})
	}

	return chunks, nil
}
