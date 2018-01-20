package dotweb

import (
	"bytes"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
)

type UploadFile struct {
	File     multipart.File
	Header   *multipart.FileHeader
	fileExt  string //file extensions
	fileName string
	fileSize int64
	content  []byte
}

func NewUploadFile(file multipart.File, header *multipart.FileHeader) *UploadFile {
	return &UploadFile{
		File:     file,
		Header:   header,
		fileName: header.Filename,
		fileExt:  filepath.Ext(header.Filename), //update for issue #99
		content:  parseFileToBytes(file),
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
func (f *UploadFile) SaveFile(fileName string) (size int, err error) {
	size = 0
	if fileName == "" {
		return size, errors.New("filename not allow empty")
	}

	fileWriter, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return size, err
	}
	defer fileWriter.Close()
	f.File.Read(f.content)
	//size, err = io.Copy(fileWriter, f.File)
	size, err = fileWriter.Write(f.content)
	return size, err
}

//get upload file extensions
func (f *UploadFile) GetFileExt() string {
	return f.fileExt
}

// Bytes returns a slice of byte hoding the UploadFile.File
func (f *UploadFile) Bytes() []byte {
	return f.content
}

func parseFileToBytes(file multipart.File) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	return buf.Bytes()
}
