package platform

import (
	"errors"
	"path/filepath"

	"github.com/FooSoft/lazarus/formats/mpq"
)

var fileState struct {
	mountPoints map[string]*mpq.Archive
	mountPaths  map[string]*mpq.Archive
}

type File struct{}

func FileMountArchive(mountPath, archivePath string) error {
	archive, err := mpq.NewFromFile(archivePath)
	if err != nil {
		return err
	}

	if fileState.mountPoints == nil {
		fileState.mountPoints = make(map[string]*mpq.Archive)
	}

	var count int
	for _, path := range archive.GetPaths() {
		resourcePath := filepath.Join(mountPath, path)
		if _, ok := fileState.mountPoints[resourcePath]; !ok {
			fileState.mountPoints[resourcePath] = archive
			count++
		}
	}

	if count == 0 {
		archive.Close()
		return errors.New("file archive could not be mounted")
	}

	return nil
}

func FileUnmountArchive(mountPath string) error {
	archive, ok := fileState.mountPoints[mountPath]
	if !ok {
		return errors.New("file archive is not mounted")
	}

	var paths []string
	for p, a := range fileState.mountPaths {
		if archive == a {
			paths = append(paths, p)
		}
	}

	for _, p := range paths {
		delete(fileState.mountPaths, p)
	}

	return nil
}

func FileUnmountAll() error {
	for _, archive := range fileState.mountPoints {
		if err := archive.Close(); err != nil {
			return err
		}
	}

	return nil
}

func FileOpen(path string) (*File, error) {
	return nil, nil
}
