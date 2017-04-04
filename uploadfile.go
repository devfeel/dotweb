package dotweb

import (
	"errors"
	files "github.com/devfeel/dotweb/framework/file"
	"io"
	"mime/multipart"
	"os"
)

type UploadFile struct {
	File     multipart.File
	Header   *multipart.FileHeader
	fileExt  string //file extensions
	fileName string
	fileSize int64
}

func NewUploadFile(file multipart.File, header *multipart.FileHeader) *UploadFile {
	return &UploadFile{
		File:     file,
		Header:   header,
		fileName: header.Filename,
		fileExt:  files.GetFileExt(header.Filename),
	}
}

// 获取文件大小的接口
type sizer interface {
	Size() int64
}

//get upload file client-local name
func (f *UploadFile) FileName() string {
	return f.fileName
}

//get upload file size
func (f *UploadFile) Size() int64 {
	if f.fileSize <= 0 {
		if sizer, ok := f.File.(sizer); ok {
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

//get upload file extensions
func (f *UploadFile) GetFileExt() string {
	return f.fileExt
}
