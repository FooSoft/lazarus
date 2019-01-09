package platform

import "github.com/FooSoft/lazarus/formats/mpq"

var fileState struct {
	mountPoints map[string]mpq.MpqArchive
}

type File struct {
}

func FileMountArchive(mountPath, archivePath string) error {
	return nil
}

func FileUnmountArchive(mountPath string) error {
	return nil
}

func FileOpen(path string) (*File, error) {
	return nil, nil
}
