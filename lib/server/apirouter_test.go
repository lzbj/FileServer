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

	/**
	req, err := http.NewRequest("GET", "http://localhost:9000/storage/download/123456/file.pdf", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler should have %s,%v,%v", rr.Code, http.StatusOK)
	}
	*/

	req1, err := http.NewRequest("POST", "http://localhost:9000/storage/upload/123456", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Errorf("handler should have %v,%v,%v", rr1.Code, http.StatusOK)
	}
}
