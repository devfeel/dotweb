package core

import (
	"net/http"
	"os"
)

//FileSystem with hide Readdir
type HideReaddirFS struct {
	FileSystem http.FileSystem
}

//File with hide Readdir
type hideReaddirFile struct {
	http.File
}

// Conforms to http.Filesystem
func (fs HideReaddirFS) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return hideReaddirFile{File: f}, nil
}

// Overrides the http.File:Readdir default implementation
func (f hideReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
