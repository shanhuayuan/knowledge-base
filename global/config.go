package global

import (
	"app/database"
	"app/object"
	"app/vectordb"
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
