package database

import (
	"errors"

	"gorm.io/gorm"
)

type GormStore struct {
	Client *gorm.DB
}

func (g *GormStore) CreateTable(table any) error {
	return g.Client.AutoMigrate(table)
}

func (g *GormStore) ListTables() ([]string, error) {
	return g.Client.Migrator().GetTables()
}

func (g *GormStore) DropTable(tableName string) error {
	return g.Client.Migrator().DropTable(tableName)
}

func (g *GormStore) HasTable(tableName string) bool {
	return g.Client.Migrator().HasTable(tableName)
}

func (g *GormStore) CreateEntry(entry any) error {
	return g.Client.Create(entry).Error
}

func (g *GormStore) DeleteEntry(entry any) error {
	return g.Client.Where(entry).Delete(entry).Error
}

func (g *GormStore) TakeEntry(entry any) error {
	return g.Client.Take(entry).Error
}

func (g *GormStore) FindEntry(query any, result any) error {
	return g.Client.Where(query).Find(result).Error
}

func (g *GormStore) IsEntryExist(entry any) (bool, error) {
	if err := g.Client.Take(entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
