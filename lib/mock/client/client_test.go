package client

import (
	"testing"
)

func TestPostFile(t *testing.T) {
	targetUrl := "http://localhost:9000/storage/upload/123456"
	filename := "file.input"
	status, err := postFile(targetUrl, filename)
	if err != nil {
		t.Fail()
	}

	if status != "200 OK" {
		t.Fail()
	}
}

func TestGetFile(t *testing.T) {
	targetUrl := "http://localhost:9000/storage/123456/file.pdf"
	status, err := getFile(targetUrl)
	if err != nil {
		t.Fail()
	}

	if status != "200 OK" {
		t.Fail()
	}
}
