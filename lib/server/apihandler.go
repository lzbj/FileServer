package server

import (
	"fmt"
	"net/http"
)

func (a storageServerHander) UploadHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %s\n", "hello")
}
