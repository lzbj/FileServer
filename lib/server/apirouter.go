package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	storageServerAPIPathPrefix = "/storage"
)

type storageServerHandler struct {
}

func RegisterStorageServerRouter(router *mux.Router) {
	storageAPI := storageServerHandler{}
	apiRouter := router.PathPrefix(storageServerAPIPathPrefix).Subrouter()
	apiRouter.Methods(http.MethodPost).HandlerFunc(storageAPI.UploadHandler)
	apiRouter.Methods(http.MethodGet).HandlerFunc(storageAPI.DownloadHandler)
}
