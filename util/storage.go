package util

import (
	"context"
	"fmt"
	"github.com/lzbj/FileServer/util/lock"
	"sync"
	"time"
)

// FSLayer implements primitives for FS Layer.
type FSStorage interface {
	CreateDir(ctx context.Context, network string, location string) error
	DeleteDir(ctx context.Context, network string, location string) error
	DeleteFile(ctx context.Context, network string, location string, fname string) error
	CreateFile(ctx context.Context, network string, location string, fname string) error
	// TODO: add more interfaces.
}

type FStorage struct {
}

func (fs *FStorage) CreateDir(ctx context.Context, network string, location string) error {
	s := fmt.Sprintf("FS Create network %s in location %s", network, location)
	fmt.Println(s)
	return nil
}

func (fs *FStorage) DeleteDir(ctx context.Context, network string, location string) error {
	s := fmt.Sprintf("FS Delete network %s in location %s", network, location)
	fmt.Println(s)
	return nil
}

func (fs *FStorage) DeleteFile(ctx context.Context, network string, location string, fname string) error {
	s := fmt.Sprintf("FS delete %s in network %s ", fname, network)
	fmt.Println(s)
	return nil
}

func (fs *FStorage) CreateFile(ctx context.Context, network string, location string, fname string) error {
	s := fmt.Sprintf("FS create %s in network %s ", fname, network)
	fmt.Println(s)
	return nil
}

type S3Storage interface {
	CreateBucket(ctx context.Context, network string, location string) error
	DeleteBucket(ctx context.Context, network string, location string) error
	DeleteObject(ctx context.Context, network string, location string, fname string) error
	CreateObject(ctx context.Context, network string, location string, fname string) error
	// TODO: add more interfaces.
}

type S3torage struct {
}

func (ss S3torage) CreateBucket(ctx context.Context, network string, location string) error {
	s := fmt.Sprintf("S3 Create network %s in location %s", network, location)
	fmt.Println(s)
	return nil
}

func (ss S3torage) DeleteBucket(ctx context.Context, network string, location string) error {
	s := fmt.Sprintf("S3 Create network %s in location %s", network, location)
	fmt.Println(s)
	return nil
}
func (ss S3torage) CreateObject(ctx context.Context, network string, location string, fname string) error {
	s := fmt.Sprintf("S3 Create network %s in location %s", network, location)
	fmt.Println(s)
	return nil
}

func (ss S3torage) DeleteObject(ctx context.Context, network string, location string, fname string) error {
	s := fmt.Sprintf("S3 Create network %s in location %s", network, location)
	fmt.Println(s)
	return nil
}

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
	fs     *FStorage
}
