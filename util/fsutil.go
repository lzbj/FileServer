package util

import "context"

// FSLayer implements primitives for FS Layer.
type FSLayer interface {
	CreateBucket(ctx context.Context, network string, location string) error
	DeleteBucket(ctx context.Context, network string, location string) error
	DeleteFile(ctx context.Context, network string, location string, fname string) error
	CreateFile(ctx context.Context, network string, location string, fname string) error
	// TODO: add more interfaces.
}
