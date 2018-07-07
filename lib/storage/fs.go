package storage

import (
	"github.com/lzbj/FileServer/util"
)

/**
func NewFSSys(fsPath string) (util.FSStorage, error) {
	fsSys := FSSys{
		fsPath: fsPath,
		rwPool: &fsIOPool{readerMap: make(map[string]*lock.RLockedFile)},
		fs: &util.FStorage{},
	}
	return fsSys, nil
}
*/

func NewFStorage(fsPath string) (util.FSStorage, error) {
	return &util.FStorage{}, nil
}
