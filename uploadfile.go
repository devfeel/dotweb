package dotweb

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

type UploadFile struct {
	File     multipart.File
	Header   *multipart.FileHeader
	fileSize int64
}

// 获取文件大小的接口
type Sizer interface {
	Size() int64
}

//get upload file client-local name
func (f *UploadFile) FileName() string {
	return f.Header.Filename
}

//get upload file size
func (f *UploadFile) Size() int64 {
	if f.fileSize <= 0 {
		if sizer, ok := f.File.(Sizer); ok {
			f.fileSize = sizer.Size()
		}
	}
	return f.fileSize
}

//save file in server-local with filename
func (f *UploadFile) SaveFile(fileName string) (size int64, err error) {
	size = 0
	if fileName == "" {
		return size, errors.New("filename not allow empty")
	}

	fileWriter, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return size, err
	}
	defer fileWriter.Close()
	size, err = io.Copy(fileWriter, f.File)
	return size, err
}
