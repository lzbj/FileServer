package util

import (
	"context"
	"testing"
)

func TestFSStorage(t *testing.T) {

	var s FSStorage
	fs := &FStorage{}
	s = fs
	ctx := context.Background()
	err := s.CreateDir(ctx, "97753", "/opt/storage")
	if err != nil {
		t.Fail()
	}

}

func TestS3torage(t *testing.T) {
	var s S3Storage
	s3 := &S3torage{}
	s = s3
	ctx := context.Background()
	err := s.CreateBucket(ctx, "97754", "/opt/tmp")
	if err != nil {
		t.Fail()
	}

}
