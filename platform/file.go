package platform

import "github.com/FooSoft/lazarus/formats/mpq"

var fileState struct {
	mountPoints map[string]*mpq.Archive
}

type File struct {
}

func FileMountArchive(mountPath, archivePath string) error {
	return nil
}

func FileUnmountArchive(mountPath string) error {
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
