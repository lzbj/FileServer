package status

import (
	"net/http"
	"github.com/lzbj/FileServer/lib/server"
	"github.com/gorilla/mux"
	"fmt"
)

func (a storageStatusHander) StatusHandler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Upload Handler: %s\n", "hello")
	vars := mux.Vars(r)

	networkid := vars["networkid"]
	filename := vars["filename"]
	query := &server.QueryRequest{
		Network:  networkid,
		Filename: filename,
	}
	fmt.Printf("received query request %s %s\n", networkid, filename)
	go queryOp(query)

	select {
	case result := <-server.GlobalOpeResChan:
		fmt.Fprintf(w, "Upload Handler :%s \n", result)
	}
}

var queryOp = func(query *server.QueryRequest) {
	server.GlobalQueryRequestChan <- query
}
