package server

import (
	"github.com/dustin/go-humanize"
	"github.com/lzbj/FileServer/util"
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
	GlobalBackEndFSSys        util.FSLayer
)

type BackendType int

const (
	globalMaxFileSize             = 1 * humanize.GiByte
	Unknown           BackendType = iota
	// File System Backend storage
	FS
	// S3 Backend storage
	S3
)
