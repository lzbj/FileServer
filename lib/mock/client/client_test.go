package client

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"
)

func TestPostFile(t *testing.T) {
	times := 5
	targetUrl := "http://127.0.0.1:9000/storage/upload/123456"
	txt := "testtxt.txt"
	pdf := "testpdf.pdf"
	png := "testpng.png"
	pptx := "testpptx.pptx"

	txtBytes, err := ioutil.ReadFile(txt)
	if err != nil {
		panic(err)
	}

	pdfBytes, err := ioutil.ReadFile(pdf)
	if err != nil {
		panic(err)
	}

	pngBytes, err := ioutil.ReadFile(png)
	if err != nil {
		panic(err)
	}

	pptBYtes, err := ioutil.ReadFile(pptx)
	if err != nil {
		panic(err)
	}
	txtBuffer := bytes.NewBuffer(txtBytes)
	pdfBuffer := bytes.NewBuffer(pdfBytes)
	pngBuffer := bytes.NewBuffer(pngBytes)
	pptBuffer := bytes.NewBuffer(pptBYtes)

	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(i int) {
			name := strconv.Itoa(i)
			status, err := postBuffer(targetUrl, name+".txt", txtBuffer)
			if err != nil {
				t.Fail()
			}

			if status != "200 OK" {
				t.Fail()
			}
			status1, err := postBuffer(targetUrl, name+".pdf", pdfBuffer)
			if err != nil {
				t.Fail()
			}

			if status1 != "200 OK" {
				t.Fail()
			}
			status2, err := postBuffer(targetUrl, name+".png", pngBuffer)
			if err != nil {
				t.Fail()
			}

			if status2 != "200 OK" {
				t.Fail()
			}
			status3, err := postBuffer(targetUrl, name+".ppt", pptBuffer)
			if err != nil {
				t.Fail()
			}

			if status3 != "200 OK" {
				t.Fail()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

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
