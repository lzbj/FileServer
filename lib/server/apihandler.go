package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/minio/minio/cmd/logger"
	"io"
	"net/http"
	"os"
)

// UploadHandler store the uploaded file current app `test` folder
func (a storageServerHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload Handler: %s\n", "hello")
	upload := vars["uploadfile"]
	fmt.Println(upload)
	r.ParseMultipartForm(32 << 1)
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

// UploadHandlerStorage store the file in the specified parameter network
// folder with the backend fs.
func (a storageServerHandler) UploadHandlerStorage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload Handler: %s\n", "hello")
	upload := vars["uploadfile"]
	fmt.Println(upload)
	r.ParseMultipartForm(32 << 1)
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
	fmt.Fprintf(w, "formfile name: %s\n", handler.Header)

}
func (a storageServerHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	networkid := vars["networkid"]
	filename := vars["filename"]
	fmt.Println(vars)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Download Handler: %s,%s,%s\n", "hello1", networkid, filename)
}
