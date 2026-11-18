package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KnowledgeBase struct {
	Id         string    `json:"id" uri:"kbId" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	ModelId    string    `json:"modelId"`
	ChunkCount uint      `json:"chunkCount" gorm:"default:0"`
	DocCount   uint      `json:"docCount" gorm:"default:0"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	deletedAt  sql.NullTime
}

func (k *KnowledgeBase) BeforeCreate(tx *gorm.DB) error {
	if k.Name == "" || k.Desc == "" || k.ModelId == "" {
		return errors.New("param error")
	}

	if k.Id == "" {
		k.Id = uuid.NewString()
	}

	return nil
}
