package server

import (
	"github.com/dustin/go-humanize"
	"github.com/lzbj/FileServer/util"
	"github.com/lzbj/FileServer/util/cache"
	"github.com/lzbj/FileServer/util/db"
	"github.com/lzbj/FileServer/util/monitor"
	"os"
	"time"
	"fmt"
	"mime/multipart"
	"context"
	"io"
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
	GlobalOperation           Operation
	GlobalOperationResult     OperationResult
	GlobalQueryRequest        QueryRequest
	GlobalOperationChannel    = make(chan *Operation)
	GlobalOpeResChan          = make(chan *OperationResult)
	GlobalQueryRequestChan    = make(chan *QueryRequest)
	results                   = make(map[string]*OperationResult)
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

type Operation struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
	Operation  string
	Content    []byte
	Network    string
}

type OperationResult struct {
	Networkid   string
	FileName    string
	Result      string
	ErrorDetail string
}

type QueryRequest struct {
	Filename string
	Network  string
}

var AsyncOps = func(stop <-chan interface{}, oc chan *Operation, rc chan *OperationResult, checkInterval time.Duration) {
	check := time.Tick(checkInterval)
	for {
		select {
		case <-stop:
			return
		case <-check:
			continue
		case op := <-oc:
			if op.Operation == "create" {
				f, err := GlobalBackEndFSSys.CreateFile(context.Background(), op.Network, op.FileHeader.Filename)
				_, err = io.Copy(f, op.File)
				if err != nil {
					//write operation fail result
					result := &OperationResult{
						Networkid:   op.Network,
						FileName:    op.FileHeader.Filename,
						Result:      "fail",
						ErrorDetail: "create file failed",
					}
					rc <- result
				} else {
					// write operation success result
					result := &OperationResult{
						Networkid:   op.Network,
						FileName:    op.FileHeader.Filename,
						Result:      "success",
						ErrorDetail: "create file success",
					}
					rc <- result
				}
			}
		default:
		}
	}
}

var QueryProducer = func(stop <-chan interface{}, opChannel chan *OperationResult) {
	for {
		select {
		case <-stop:
			return
		case result := <-opChannel:
			//update the result map
			name := result.FileName
			fmt.Printf("put the operation result of file %s into result map\n", name)
			results[name] = result
		default:
		}
	}
}

var QueryConsumer = func(stop <-chan interface{}, queryRequestChannel chan *QueryRequest) {
	for {
		select {
		case <-stop:
			return
		case query := <-queryRequestChannel:
			fmt.Printf("received query %s\n", query.Filename)
			_, ok := results[query.Filename]
			fmt.Printf("handle query result %v\n", ok)
			if ok {
				GlobalOpeResChan <- results[query.Filename]
			} else {
				errorResult := &OperationResult{
					Networkid:   "",
					FileName:    "",
					Result:      "fail",
					ErrorDetail: "operation failed",
				}
				GlobalOpeResChan <- errorResult
			}
		default:

		}
	}
}

