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

	apiRouter := router.PathPrefix(storageServerAPIPathPrefix).Path("/{networkid:[0-9]+}").Subrouter()
	apiRouter = router.PathPrefix(storageServerAPIPathPrefix).Path("/{networkid:[0-9]+}/{filename}").Subrouter()
	apiRouter.Methods(http.MethodPost).HandlerFunc(storageAPI.UploadHandler)
	apiRouter.Methods(http.MethodGet).HandlerFunc(storageAPI.DownloadHandler)
}
