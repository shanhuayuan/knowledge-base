package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelType string

const (
	ModelTypeChat  ModelType = "chat"
	ModelTypeEmbed ModelType = "embed"
)

type Provider struct {
	Id        string    `json:"id" uri:"providerId" gorm:"primaryKey"`
	Provider  string    `json:"provider"`
	ApiKey    string    `json:"apiKey"`
	BaseUrl   string    `json:"baseUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	deletedAt sql.NullTime
}

func (p *Provider) BeforeCreate(tx *gorm.DB) error {
	if p.Provider == "" || p.ApiKey == "" || p.BaseUrl == "" {
		return errors.New("param error")
	}

	if p.Id == "" {
		p.Id = uuid.NewString()
	}

	return nil
}

type Model struct {
	Id         string    `json:"id" uri:"modelId" gorm:"primaryKey"`
	Name       string    `json:"name"`
	ProviderId string    `json:"providerId"`
	Type       ModelType `json:"type"`
	Dim        int64     `json:"dim"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	deletedAt  sql.NullTime
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.Name == "" || m.ProviderId == "" || m.Type == "" || m.Dim == 0 {
		return errors.New("param error")
	}

	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	return nil
}
