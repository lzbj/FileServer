package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	APIPathPrefix = "/storage"
)

type storageServerHandler struct {
}

func RegisterStorageServerRouter(router *mux.Router) {
	storageAPI := storageServerHandler{}

	apiRouter := router.PathPrefix(APIPathPrefix).Path("/upload/{networkid:[0-9]+}").Subrouter()
	//apiRouter = router.PathPrefix(APIPathPrefix).Path("/download/{networkid:[0-9]+}/{filename}").Subrouter()
	//apiRouter.Methods(http.MethodPost).HandlerFunc(storageAPI.UploadHandler)
	apiRouter.Methods(http.MethodPost).HandlerFunc(storageAPI.UploadHandlerStorage)
	//apiRouter.Methods(http.MethodGet).HandlerFunc(storageAPI.DownloadHandler)
}

func RegisterStorageServerRouterDownload(router *mux.Router) {
	storageAPI := storageServerHandler{}

	apiRouter := router.PathPrefix(APIPathPrefix).Path("/download/{networkid:[0-9]+}/{filename}").Subrouter()
	apiRouter.Methods(http.MethodGet).HandlerFunc(storageAPI.DownloadHandler)
}
