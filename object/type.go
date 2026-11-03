package object

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type (
	ObjectInfo struct {
		// Id         string
		Name        string `json:"name"`
		Size        int64  `json:"size"`
		ContentType string `json:"contentType"`
		Error       bool   `json:"error"`
	}

	ObjectStore interface {
		ListObjects(ctx context.Context, option ListObjectsOption) ListObjectsResult
		UploadObject(ctx context.Context, option UploadObjectOption) error
		GetObjectInfo(ctx context.Context, option GetObjectOption) (*ObjectInfo, error)
		DeleteObject(ctx context.Context, option DeleteObjectOption) error
		DeleteObjects(ctx context.Context, option DeleteObjectOption)
	}

	ListObjectsOption struct {
		BucketName string `form:"bucketName"`
		Prefix     string `form:"prefix"`
	}

	ListObjectsResult struct {
		ObjectsInfo []ObjectInfo `json:"objectsInfo"`
		Total       int          `json:"total"`
	}

	DeleteObjectOption struct {
		BucketName string `form:"bucketName"`
		Prefix     string `form:"prefix"`
	}

	DeleteObjectsResult struct {
		Objects []ObjectInfo
		Total   int
	}

	GetObjectOption struct {
		BucketName string `form:"bucketName"`
		Prefix     string `form:"prefix"`
	}

	UploadObjectOption struct {
		BucketName  string `form:"bucketName"`
		Prefox      string `form:"prefix"`
		File        multipart.File
		Size        int64
		ContentType string
	}

	MinioStore struct {
		Client *minio.Client
	}
)

func (r *ListObjectsResult) Append(info ObjectInfo) {
	r.ObjectsInfo = append(r.ObjectsInfo, info)
	r.Total += 1
}
