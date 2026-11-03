package knowledge

import "io"

type (
	UploadFileOption struct {
		FileName string
		File     io.Reader
		FileSize int64
		FileType string
		Target   string
	}

	FileStore struct {
		//FileStore的Id由KnowledgeStore的Id决定
		Id         string
		BucketName string
	}
)

type MinioFileStore struct {
	Id         string
	BucketName string
}

func (s *MinioFileStore) UploadFile(option *UploadFileOption) {

}
