package storage

import (
	"context"
	"github.com/lzbj/FileServer/util"
	"github.com/lzbj/FileServer/util/lock"
	"sync"
	"time"
)

type BackendType int

const (
	Unknown BackendType = iota
	// File System Backend storage
	FS
	// S3 Backend storage
	S3
)

// StorageInfo - represents total capacity of underlying storage.
type StorageInfo struct {
	Used uint64

	Backend struct {
		Type BackendType
	}
}

// BucketInfo - represents bucket metadata.
type BucketInfo struct {
	Name    string
	Created time.Time
}

// FileInfo - represents file metadata.
type FileInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the file.
	Name string
	// Data and time when the object was last modified.
	ModTime time.Time

	// Total object Size.
	Size int64

	//IsDir
	IsDir bool

	// ContentType
	ContentType string

	//ContentEncoding
	ContentEncoding string
}

type fsIOPool struct {
	sync.Mutex
	readerMap map[string]*lock.RLockedFile
}

// BucketInfo

type FSSys struct {
	// File System Path
	fsPath string
	rwPool *fsIOPool
}

func NewFSSys(fsPath string) (util.FSLayer, error) {
	fsSys := &FSSys{
		fsPath: fsPath,
		rwPool: &fsIOPool{readerMap: make(map[string]*lock.RLockedFile)},
	}
	return fsSys, nil
}

func (fs *FSSys) Shutdown(ctx context.Context) error {
	return nil
}

func (fs *FSSys) CreateBucket(ctx context.Context, network string, location string) error {
	return nil
}

func (fs *FSSys) DeleteBucket(ctx context.Context, network string, location string) error {
	return nil
}

func (fs *FSSys) DeleteFile(ctx context.Context, network string, location string, fname string) error {
	return nil
}

func (fs *FSSys) CreateFile(ctx context.Context, network string, location string, fname string) error {
	return nil
}
