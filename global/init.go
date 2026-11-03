package global

import (
	"app/object"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type Runtime struct {
	Config *Config
	Store  object.ObjectStore
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

	OverallRuntime.Store = &object.MinioStore{Client: client}
	return nil
}
