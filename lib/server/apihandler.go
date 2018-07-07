package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

func (a storageServerHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload Handler: %s\n", "hello")
	_ = vars["uploadfile"]
	r.ParseMultipartForm(32 << 1)
	file, handler, err := r.FormFile("uploadfile")
	network := r.FormValue("network")
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
	err = GlobalBackEndFSSys.CreateDir(context.Background(), network, "opt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	io.Copy(f, file)

	fmt.Fprintf(w, "formfile name: %s\n", handler.Header)

}

func (a storageServerHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Download Handler: %s\n", "hello")
}
