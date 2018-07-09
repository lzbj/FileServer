package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterStorageServerRouter(t *testing.T) {
	router := mux.NewRouter().SkipClean(true)

	//Register upload api router

	RegisterStorageServerRouter(router)

	req, err := http.NewRequest("GET", "http://localhost:9000/storage/123456/file.pdf", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler should have %s,%v,%v", rr.Code, http.StatusOK)
	}
}
