package status

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	storageServerStatusPrefix = "/status"
)

type storageStatusHander struct {
}

func RegisteStatusRouter(router *mux.Router) {
	statusAPI := storageStatusHander{}

	apiRouter := router.PathPrefix(storageServerStatusPrefix).Subrouter()

	apiRouter.Methods(http.MethodGet).Path("/query/{networkid}/{filename}").HandlerFunc(statusAPI.StatusHandler)
}
