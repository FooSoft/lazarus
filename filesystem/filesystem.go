package filesystem

import (
	"bytes"
	"io"

	"github.com/icza/mpq"
)

type FileSystem interface {
	Mount(root, path string) error
	List() ([]string, error)
	Open(path string) (io.ReadSeeker, error)
}

func New() FileSystem {
	return new(filesystem)
}

type filesystem struct {
	db *mpq.MPQ
}

func (fs *filesystem) Mount(root, path string) error {
	db, err := mpq.NewFromFile(path)
	if err != nil {
		return err
	}

	fs.db = db
	return nil
}

func (fs *filesystem) List() ([]string, error) {
	data, err := fs.db.FileByName("(listfile)")
	if err != nil {
		return nil, err
	}

	return []string{string(data)}, nil
}

func (fs *filesystem) Open(path string) (io.ReadSeeker, error) {
	data, err := fs.db.FileByName(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
