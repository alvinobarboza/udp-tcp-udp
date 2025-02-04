package filehandler

import (
	"errors"
	"os"
)

type FileHandler interface {
	Write([]byte) ([]byte, error)
	NewFile(string) error
}

func NewFileHandler() FileHandler {
	return &fileHandler{}
}

type fileHandler struct {
	fileDesc *os.File
}

func (f *fileHandler) Write(data []byte) ([]byte, error) {

	if f.fileDesc == nil {
		return nil, ErrNoFileAvailable
	}

	_, errW := f.fileDesc.Write(data)
	if errW != nil {
		return nil, errW
	}
	return nil, nil
}

func (f *fileHandler) NewFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.fileDesc = file

	return nil
}

func (f *fileHandler) Close() error {
	err := f.fileDesc.Close()

	f.fileDesc = nil

	return err
}

var ErrNoFileAvailable = errors.New("no file created! try calling NewFile")
