package util

import (
	"context"
	"fmt"
	"github.com/lzbj/FileServer/util/lock"
	"github.com/minio/minio/cmd/logger"
	"github.com/spf13/afero"
	"os"
	"sync"
	"time"
)

// FSLayer implements primitives for FS Layer.
type FSStorage interface {
	CreateDir(ctx context.Context, network string) error
	CreateDirNew(ctx context.Context, network string) error
	DeleteDir(ctx context.Context, network string) error
	DeleteFile(ctx context.Context, network string, fname string) error
	CreateFile(ctx context.Context, network string, fname string) (afero.File, error)
	GetFile(ctx context.Context, network string, fname string) (afero.File, error)
	// TODO: add more interfaces.
}

func NewFStorage(path string) (*FStorage, error) {
	s := &FStorage{}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.err = make(chan error)
	s.report = make(chan StorageStatus)

	fs := afero.NewOsFs()

	exist, err := afero.Exists(fs, path)
	if err != nil {
		logger.Info("error %s", err)
		return nil, err
	}
	if !exist {
		err = fs.Mkdir(path, 0755)
		if err != nil {
			logger.Info("error %s", err)
			return nil, err
		}
	}
	s.fs = afero.NewBasePathFs(fs, path)
	return s, nil
}

// FStorage implements the FSStorage interfaces.
type FStorage struct {
	// lock for the storage
	lock sync.Mutex
	// The root dir store the files
	fs afero.Fs
	// Error chan for errors in the storage.
	err chan error
	// The ticker for the FS usage check intervals.
	ticker time.Ticker
	report chan StorageStatus
}

// StorageStatus represents the FS storage status.
type StorageStatus struct {
	diskSize uint32
	usage    float32
}

func (fs *FStorage) CreateDir(ctx context.Context, network string) error {
	s := fmt.Sprintf("FS Create network %s in location %s", network)
	fmt.Println(s)
	return nil
}

func (fs *FStorage) CreateDirNew(ctx context.Context, network string) error {
	return fs.fs.Mkdir(network, 0755)
}

func (fs *FStorage) DeleteDir(ctx context.Context, network string) error {
	s := fmt.Sprintf("FS Delete network %s in location %s", network)
	fmt.Println(s)
	return nil
}

func (fs *FStorage) DeleteFile(ctx context.Context, network string, fname string) error {
	s := fmt.Sprintf("FS delete %s in network %s ", fname, network)
	fmt.Println(s)
	return nil
}

func (fs *FStorage) CreateFile(ctx context.Context, network string, fname string) (afero.File, error) {
	s := fmt.Sprintf("FS create %s in network %s ", fname, network)
	fmt.Println(s)
	f, err := fs.fs.Create(network + "/" + fname)
	if err != nil {
		logger.Info("error happened during create file %s", err)
		return nil, err
	}
	return f, nil
}

func (fs *FStorage) GetFile(ctx context.Context, network string, fname string) (afero.File, error) {
	wFileName := network + "/" + fname
	f, err := fs.fs.OpenFile(wFileName, os.O_RDONLY, 0666)
	if err != nil {
		logger.Info("error happened during get file %s", err)
		return nil, err
	}
	return f, nil
}

func (fs *FStorage) CreateFileNew(ctx context.Context, network string, fname string) error {
	_, err := fs.fs.Create(network)
	if err != nil {
		return err
	}
	//fs.fs.Create()
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
