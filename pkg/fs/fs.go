package fs

import (
	"mime/multipart"
)

type Uploader interface {
	Init([]byte) error
	UploadHTTP(mh *multipart.FileHeader) (*FileInfo, error)
	UploadLocal(dst string) (*FileInfo, error)
}

type FileInfo struct {
	FullURL  string `json:"full_url"`
	BaseURL  string `json:"base_url"`
	Domain   string `json:"domain"`
	FileSize int64  `json:"file_size"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
}
