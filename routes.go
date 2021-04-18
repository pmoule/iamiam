package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pmoule/iamiam/iamiam"
)

// InitRouter initialises all all REST routes.
func InitRouter(validUserInfos []*iamiam.UserInfo, template string) *mux.Router {
	knownUsers = append(knownUsers, validUserInfos...)
	loginTemplate = template

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/auth", auth).Methods(http.MethodGet)
	r.HandleFunc("/use", useEmail).Methods(http.MethodPost)
	r.HandleFunc("/token", token).Methods(http.MethodPost)
	r.HandleFunc("/info", info).Methods(http.MethodGet)

	webDirName := ""
	staticContentHandler := staticContent("./" + webDirName)
	r.PathPrefix("/").Handler(staticContentHandler)

	return r
}

func staticContent(dirName string) http.Handler {
	staticContentDir := http.Dir(dirName)
	return http.FileServer(staticContentDir)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
