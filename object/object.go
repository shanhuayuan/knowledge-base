package object

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (s *MinioStore) ListObjects(ctx context.Context, option ListObjectsOption) ListObjectsResult {
	result := ListObjectsResult{}

	listObjectsOption := minio.ListObjectsOptions{}
	for object := range s.Client.ListObjects(ctx, option.BucketName, listObjectsOption) {
		info := ObjectInfo{
			Name:        object.Key,
			Size:        object.Size,
			ContentType: object.ContentType,
			Error:       object.Err != nil,
		}
		result.Append(info)
	}
	return result
}

func (s *MinioStore) GetObjectInfo(ctx context.Context, option GetObjectOption) (*ObjectInfo, error) {
	info, err := s.Client.StatObject(ctx, option.BucketName, option.Prefix, minio.StatObjectOptions{})

	if err != nil {
		return nil, err
	}

	result := &ObjectInfo{
		Name:        info.Key,
		Size:        info.Size,
		ContentType: info.ContentType,
		Error:       info.Err != nil,
	}

	return result, nil
}

func (s *MinioStore) DeleteObject(ctx context.Context, option DeleteObjectOption) error {
	deleteObjectOption := minio.RemoveObjectOptions{}

	if err := s.Client.RemoveObject(ctx, option.BucketName, option.Prefix, deleteObjectOption); err != nil {
		return err
	}

	return nil
}

func (s *MinioStore) DeleteObjects(ctx context.Context, option DeleteObjectOption) {
	deleted := make([]string, 0)

	objectCh := make(chan minio.ObjectInfo)
	removeCh := minio.RemoveObjectsOptions{}

	go func() {
		defer close(objectCh)
		listObjectsOption := minio.ListObjectsOptions{Prefix: option.Prefix, Recursive: true}

		for object := range s.Client.ListObjects(ctx, option.BucketName, listObjectsOption) {
			deleted = append(deleted, object.Key)
			objectCh <- object
		}
	}()

	s.Client.RemoveObjects(ctx, option.BucketName, objectCh, removeCh)
}

func (s *MinioStore) UploadObject(ctx context.Context, option UploadObjectOption) error {
	putObjectOption := minio.PutObjectOptions{ContentType: option.ContentType}
	_, err := s.Client.PutObject(ctx, option.BucketName, option.Prefox, option.File, option.Size, putObjectOption)
	return err
}
