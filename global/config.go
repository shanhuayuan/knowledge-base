package global

import (
	"github.com/shanhuayuan/database"
	"github.com/shanhuayuan/object"
	"github.com/shanhuayuan/vectordb"
)

type Config struct {
	App struct {
		Port string
	}

	Knowledge struct {
		Milvus vectordb.StoreConfig
	}

	Database struct {
		Mysql database.StoreConfig
	}

	Object struct {
		Minio object.StoreConfig
	}
}
