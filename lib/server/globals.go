package server

import (
	"github.com/dustin/go-humanize"
	"os"
)

var (
	GlobalPort                = "127.0.0.1:9000"
	GlobalFSPath              = "/opt/fileStorageServer"
	GlobalCacheFSPath         = "/opt/cacheFileStorageServer"
	StorageDefaultBackEndType = "FS"
	StorageS3BackendType      = "AWSS3"
	GlobalHTTPServerErrorCh   = make(chan error)
	GlobalOSSignalCh          = make(chan os.Signal, 1)
)

const (
	globalMaxFileSize = 1 * humanize.GiByte
)
