package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"testTaskIdeasoft/repository"
)

func NewHandler(cmds chan repository.Command) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/request", RecoveryMiddleware(Request(cmds))).Methods("GET")
	router.HandleFunc("/admin/requests", RecoveryMiddleware(GetAnalytics(cmds))).Methods("GET")
	return router
}

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	}
}

func Request(cmds chan repository.Command) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ch := make(chan string)
		cmds <- repository.Command{
			Typ:   repository.View,
			ReqCh: ch,
		}
		resp := <-ch
		resp = fmt.Sprintf("%v\n", resp)
		_, err := writer.Write([]byte(resp))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetAnalytics(cmds chan repository.Command) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ch := make(chan map[string]int)
		cmds <- repository.Command{
			Typ:        repository.GetViews,
			AnalyticCh: ch,
		}
		resp := <-ch
		for k, v := range resp {
			_, err := writer.Write([]byte(fmt.Sprintf("%v - %v\n", k, v)))
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
