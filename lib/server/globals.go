package server

import (
	"github.com/dustin/go-humanize"
	"github.com/lzbj/FileServer/util"
	"github.com/lzbj/FileServer/util/cache"
	"github.com/lzbj/FileServer/util/db"
	"github.com/lzbj/FileServer/util/monitor"
	"os"
)

var (
	GlobalPort                = "127.0.0.1:9000"
	GlobalFSPath              = "/opt/FileStorage"
	GlobalFSInitalized        = false
	GlobalCacheFSPath         = "/opt/CacheFileStorage"
	StorageDefaultBackEndType = "FS"
	StorageS3BackendType      = "AWSS3"
	GlobalHTTPServerErrorCh   = make(chan error)
	GlobalOSSignalCh          = make(chan os.Signal, 1)
	GlobalBackEndFSSys        util.FSStorage
	GlobalDBCache             cache.RedisCache
	GlolabDBConn              db.DBConn
	GlobalMonitorSys          monitor.MonitorSys
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

//Errors
