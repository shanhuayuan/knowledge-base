package database

import (
	"fmt"
)

type StoreConfig struct {
	User     string
	Password string
	Port     string
	Database string
}

type DropTableOption struct {
	TableName string `uri:"tableName"`
}

func (c *StoreConfig) ToDsn() string {
	return fmt.Sprintf("%s:%s@tcp(mysql:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Port, c.Database)
}

type Store interface {
	CreateTable(table any) error
	DropTable(string) error
	ListTables() ([]string, error)
	HasTable(string) bool
	CreateEntry(entry any) error
	DeleteEntry(entry any) error
	TakeEntry(entry any) error
	FindEntry(query any, result any) error
	IsEntryExist(entry any) (bool, error)
}
