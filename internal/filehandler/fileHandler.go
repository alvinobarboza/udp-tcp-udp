package filehandler

import (
	"errors"
	"os"
)

type FileHandler interface {
	Write([]byte, chan error)
	NewFile(string) error
	CloseConn()
}

func NewFileHandler() FileHandler {
	return &fileHandler{}
}

type fileHandler struct {
	fileDesc *os.File
}

func (f *fileHandler) Write(data []byte, err chan error) {

	if f.fileDesc == nil {
		err <- ErrNoFileAvailable
		return
	}

	_, errW := f.fileDesc.Write(data)
	if errW != nil {
		err <- errW
		return
	}
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

func (f *fileHandler) CloseConn() {}

var ErrNoFileAvailable = errors.New("no file created! try calling NewFile")
