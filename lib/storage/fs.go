package storage

import (
	"github.com/lzbj/FileServer/util"
	"github.com/minio/minio/cmd/logger"
)

func NewFStorage(fsPath string) (util.FSStorage, error) {
	storage, err := util.NewFStorage(fsPath)
	if err != nil {
		logger.Info("error happened during new fs storage %s", err)
		return nil, err
	}

	return storage, nil

}
