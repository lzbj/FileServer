package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	APIPathPrefix = "/storage"
	AsyncAPIPathPrefix = "/storage/async"
)

type storageServerHandler struct {
}

func RegisterStorageServerRouter(router *mux.Router) {
	storageAPI := storageServerHandler{}
	apiRouter := router.PathPrefix(APIPathPrefix).Path("/upload").Subrouter()
	apiRouter.Methods(http.MethodPost).HandlerFunc(storageAPI.UploadHandlerStorage)
}

func RegisterStorageServerRouterDownload(router *mux.Router) {
	storageAPI := storageServerHandler{}
	apiRouter := router.PathPrefix(APIPathPrefix).Path("/download/{networkid:[0-9]+}/{filename}").Subrouter()
	apiRouter.Methods(http.MethodGet).HandlerFunc(storageAPI.DownloadHandler)
}


func RegisterStorageServerRouterAsync(router *mux.Router){
	storageAPI := storageServerHandler{}
	apiRouter := router.PathPrefix(AsyncAPIPathPrefix).Path("/upload").Subrouter()
	apiRouter.Methods(http.MethodPost).HandlerFunc(storageAPI.AsyncUploadHandlerStorage)
}