package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
	"github.com/minio/minio/cmd/logger"
	"io"
	"net/http"
	"os"
)

// UploadHandlerStorage store the file in the specified parameter network
// folder with the backend fs.

type uploadResult struct {
	Result       string `json:"result"`
	DownloadLink string `json:"downloadLink"`
}

func (a storageServerHandler) UploadHandlerStorage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Upload Handler: %s\n", "hello")
	upload := vars["uploadfile"]
	fmt.Println(upload)
	r.Body = http.MaxBytesReader(w, r.Body, humanize.GByte*3)
	r.ParseForm()
	//r.ParseMultipartForm(64 << 1)
	file, handler, err := r.FormFile("uploadfile")
	network := r.FormValue("network")
	fmt.Println(vars)
	if len(network) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	logger.Info("network %s", network)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = GlobalBackEndFSSys.CreateDirNew(context.Background(), network)
	f, err := GlobalBackEndFSSys.CreateFile(context.Background(), network, handler.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	io.Copy(f, file)
	//fmt.Fprintf(w, "formfile name: %s\n", handler.Header)

	//TODO: implement upload and download URL in json format.
	//
	result := uploadResult{
		Result:       "ok",
		DownloadLink: "http://127.0.0.1:9000/storage/download" + "/" + network + "/" + handler.Filename,
	}
	s := fmt.Sprintf("%v", result)
	fmt.Println(s)
	b, _ := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s\n", b)

}
func (a storageServerHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	networkid := vars["networkid"]
	filename := vars["filename"]
	f, err := GlobalBackEndFSSys.GetFile(context.Background(), networkid, filename)
	defer f.Close()
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
	}

	stats, err := f.Stat()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}

	writer := w.(io.Writer)
	buf := make([]byte, stats.Size())
	n := 0
	for err == nil {
		n, err = f.Read(buf)
		writer.Write(buf[0:n])
	}

	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Download Handler: %s,%s,%s\n", "hello1", networkid, filename)
}

// UploadHandler store the uploaded file current app `test` folder
func (a storageServerHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload Handler: %s\n", "hello")
	upload := vars["uploadfile"]
	fmt.Println(upload)

	http.MaxBytesReader(w, r.Body, humanize.GByte*3)
	r.ParseForm()
	file, handler, err := r.FormFile("uploadfile")
	network := r.FormValue("network")
	fmt.Println(vars)
	if len(network) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	err = GlobalBackEndFSSys.CreateDir(context.Background(), network)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	io.Copy(f, file)

	fmt.Fprintf(w, "formfile name: %s\n", handler.Header)

}
