package client

import (
	"bytes"
	"fmt"
	"github.com/minio/minio/cmd/logger"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func getFile(targetUrl string) (string, error) {
	resp, err := http.Get(targetUrl)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	logger.Info(resp.Status)
	logger.Info("FileUploadServer Response: " + string(respBody))

	return resp.Status, nil
}

func postFile(targetUrl string, filename string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.CreateFormField("")
	bodyWriter.WriteField("network", "97753")
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return "", err
	}
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return "", err
	}
	defer fh.Close()
	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	logger.Info(resp.Status)
	logger.Info("FileUploadServer Response: " + string(respBody))

	return resp.Status, nil

}

/**
func uploadFile(targetUrl,filename string) {
	c := http.Client{}
	f, err := os.Open(filename)
	if err != nil {
		logger.Info("Error in opening file", err)
		return
	}
	defer f.Close()

	body := &bytes.Buffer{}

	_, err = io.Copy(body, f)
	if err != nil {
		logger.Info("Copy failed.", err)
		return
	}

	req, err := http.NewRequest("POST", targetUrl, bytes.NewReader(body.Bytes()))
	if err != nil {
		logger.Info("new request failed.", err)
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.Do(req)
	if err != nil {
		logger.Info("do request failed.", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info("error read body", err)
	}

	logger.Info(resp.Status)
	logger.Info("FileUploadServer Response: " + string(content))

}
*/
