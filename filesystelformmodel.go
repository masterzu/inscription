package main

import (
	// "encoding/json"
	// "fmt"
	"io"
)

// The storage datas using files
type FileSystemFormModel struct {
	data io.ReadSeeker
}

// Read
func (f *FileSystemFormModel) Read() (FormModel, error) {
	f.data.Seek(0, 0)
	model, err := ReadFormModel(f.data)
	return model, err
}
