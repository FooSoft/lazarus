package filesystem

import (
	"io"
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
}

func (fs *filesystem) Mount(root, path string) error {
	return nil
}

func (fs *filesystem) List() ([]string, error) {
	return nil, nil
}

func (fs *filesystem) Open(path string) (io.ReadSeeker, error) {
	return nil, nil
}
