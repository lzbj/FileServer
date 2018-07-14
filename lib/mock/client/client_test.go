package client

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"sync"
	"testing"
)

func TestPostFile(t *testing.T) {
	times := 1
	targetUrl := "http://127.0.0.1:9000/storage/upload"
	networkname := "97753"
	txt := "testtxt.txt"
	//rmvb := "interstella.rmvb"
	pdf := "testpdf.pdf"
	png := "testpng.png"
	pptx := "testpptx.pptx"

	var wg sync.WaitGroup

	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(i int) {
			txtBytes, err := ioutil.ReadFile(txt)
			if err != nil {
				panic(err)
			}

			/**
			rmvbBytes, err := ioutil.ReadFile(rmvb)
			if err != nil {
				panic(err)
			}
			*/
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
			//rmvbBuffer := bytes.NewBuffer(rmvbBytes)
			name := strconv.Itoa(i)
			status, err := postBuffer(targetUrl, networkname, name+".txt", txtBuffer)
			if err != nil {
				t.Fail()
			}

			if status != "200 OK" {
				t.Fail()
			}
			status1, err := postBuffer(targetUrl, networkname, name+".pdf", pdfBuffer)
			if err != nil {
				t.Fail()
			}

			if status1 != "200 OK" {
				t.Fail()
			}
			status2, err := postBuffer(targetUrl, networkname, name+".png", pngBuffer)
			if err != nil {
				t.Fail()
			}

			if status2 != "200 OK" {
				t.Fail()
			}
			status3, err := postBuffer(targetUrl, networkname, name+".ppt", pptBuffer)
			if err != nil {
				t.Fail()
			}

			if status3 != "200 OK" {
				t.Fail()
			}

			/**
			status4, err := postBuffer(targetUrl, networkname, name+".rmbv", rmvbBuffer)
			if err != nil {
				t.Fail()
			}

			if status4 != "200 OK" {
				t.Fail()
			}
			*/

			wg.Done()
		}(i)
	}
	wg.Wait()

}

func TestGetFile(t *testing.T) {
	targetUrl := "http://127.0.0.1:9000/storage/download/97753/0.pdf"
	status, err := getFile(targetUrl)
	if err != nil {
		t.Fail()
	}

	if status != "200 OK" {
		t.Fail()
	}
}
