package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	storageServerAPIPathPrefix = "/storage"
)

type storageServerHander struct {
}

func RegisterStorageServerRouter(router *mux.Router) {
	storageAPI := storageServerHander{}

	apiRouter := router.PathPrefix(storageServerAPIPathPrefix).Subrouter()

	apiRouter.Methods(http.MethodGet).Path("").HandlerFunc(storageAPI.UploadHandler)
}
