package stats

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	storageServerStatsPrefix = "/stats"
)

type storageStatsHander struct {
}

func RegisteStatusRouter(router *mux.Router) {
	statsAPI := storageStatsHander{}

	apiRouter := router.PathPrefix(storageServerStatsPrefix).Subrouter()

	apiRouter.Methods(http.MethodGet).Path("").HandlerFunc(statsAPI.StatsHandler)
}
