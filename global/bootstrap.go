package global

import (
	"app/chunk"
	"app/database"
	"app/object"
	"app/vectordb"
	"context"
	"errors"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Runtime struct {
	Config      *Config
	VectorStore vectordb.Store
	ObjectStore object.Store
	DataStore   database.Store
	Chunker     *chunk.Chunker
}

var OverallRuntime = Runtime{}

func InitConfig() error {
	config := new(Config)

	viper.SetConfigFile("./global/config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	OverallRuntime.Config = config
	return nil
}

func InitMinio() error {
	minioConfig := OverallRuntime.Config.Object.Minio

	options := minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, minioConfig.Token),
		Secure: minioConfig.UseSsl,
	}

	client, err := minio.New(minioConfig.Endpoint, &options)

	if err != nil {
		return err
	}

	OverallRuntime.ObjectStore = &object.MinioStore{Client: client}
	return nil
}

func InitMysql() error {
	config := OverallRuntime.Config.Database.Mysql

	db, err := gorm.Open(mysql.Open(config.ToDsn()), &gorm.Config{})

	if err != nil {
		return err
	}

	store := database.GormStore{Client: db}

	OverallRuntime.DataStore = &store
	return nil
}

func InitMilvus() error {
	milvusConfig := OverallRuntime.Config.Knowledge.Milvus

	config := &milvusclient.ClientConfig{
		Address: milvusConfig.Address + ":" + milvusConfig.Port,
	}

	client, err := milvusclient.New(context.Background(), config)

	if err != nil {
		return nil
	}

	if client == nil {
		return errors.New("client is nil")
	}

	store := vectordb.MilvusStore{Client: client}

	OverallRuntime.VectorStore = &store
	return nil
}

func InitChunker() {
	OverallRuntime.Chunker = &chunk.Chunker{}
}
