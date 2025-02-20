package filehandler

import (
	"errors"
	"os"
)

type FileHandler interface {
	Write([]byte) error
	NewFile(string) error
	CloseConn()
}

func NewFileHandler() FileHandler {
	return &fileHandler{}
}

type fileHandler struct {
	fileDesc *os.File
}

func (f *fileHandler) Write(data []byte) error {

	if f.fileDesc == nil {
		return ErrNoFileAvailable
	}

	_, errW := f.fileDesc.Write(data)
	if errW != nil {
		return errW
	}
	return nil
}

func (f *fileHandler) NewFile(filename string) error {
	file, err := os.Create(filename)
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

func (f *fileHandler) CloseConn() {}

var ErrNoFileAvailable = errors.New("no file created! try calling NewFile")
