package forward

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"

	"github.com/timberio/timber-go/logging"
)

type FileForwarder struct {
	Filename string
	Logger   logging.Logger

	file *os.File
}

func NewFileForwarder(filename string, logger logging.Logger) (*FileForwarder, error) {
	if logger == nil {
		logger = logging.DefaultLogger
	}

	file, err := openFile(filename)
	if err != nil {
		return nil, err
	}

	return &FileForwarder{
		Filename: filename,

		file: file,
	}, nil
}

func openFile(filename string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return nil, err
	}

	var file *os.File

	if _, err := os.Stat(filename); err == nil {
		file, err = os.OpenFile(filename, os.O_RDWR, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Create(filename)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func (f *FileForwarder) Forward(buffer *bytes.Buffer) {
	writer := bufio.NewWriter(f.file)

	_, err := writer.Write(buffer.Bytes())
	if err != nil {
		f.Logger.Print(err)
		return
	}

	err = writer.Flush()
	if err != nil {
		f.Logger.Print(err)
		return
	}
}
