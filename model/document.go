package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentStatus string

const (
	DocumentNotStarted = "notStarted"
	DocumentQueued     = "queued"
	DocumentProcessing = "processing"
	DocumentFinished   = "finished"
	DocumentFailed     = "failed"
)

type Document struct {
	Id         string         `json:"id" uri:"docId" gorm:"primaryKey"`
	KbId       string         `json:"kbId" uri:"kbId"`
	Name       string         `json:"name"`
	FileType   string         `json:"fileType"`
	FileSize   int64          `json:"fileSize"`
	ObjectKey  string         `json:"objectKey"`
	Status     DocumentStatus `json:"status"`
	ChunkCount int            `json:"chunkCount"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	deletedAt  sql.NullTime
}

func (d *Document) BeforeCreate(_ *gorm.DB) error {
	if d.Id == "" {
		d.Id = uuid.NewString()
	}

	return nil
}

type ContentPart struct {
	Id         string    `json:"id" uri:"contentPartId" gorm:"primaryKey"`
	DocumentId string    `json:"documentId"`
	Type       string    `json:"type"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	deletedAt  sql.NullTime
}

func (c *ContentPart) BeforeCreate(_ *gorm.DB) error {
	if c.Id == "" {
		c.Id = uuid.NewString()
	}

	return nil
}
