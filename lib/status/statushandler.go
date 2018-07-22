package status

import (
	"net/http"
	"github.com/lzbj/FileServer/lib/server"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
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
		s := fmt.Sprintf("%v", result)
		fmt.Println(s)
		b, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%s \n", b)
	}
}

var queryOp = func(query *server.QueryRequest) {
	server.GlobalQueryRequestChan <- query
}
