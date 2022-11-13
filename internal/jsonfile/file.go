package jsonfile

import (
	"bufio"
	"encoding/json"
	"fileconverter/internal/model"
	"os"
)

type File struct {
	*os.File
}

func New(f *os.File) *File {
	return &File{f}
}

type data struct {
	Objects model.Objects `json:"objects"`
}

func (file *File) ToObjects() (model.Objects, error) {
	var d data

	decoder := json.NewDecoder(bufio.NewReader(file.File))

	err := decoder.Decode(&d)
	if err != nil {
		return model.Objects{}, err
	}

	return d.Objects, nil
}
