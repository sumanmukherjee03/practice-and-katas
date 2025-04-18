package app

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sumanmukherjee03/practice-and-katas/go/rest-api-project/bookstore_items-api/clients/elasticsearch"
)

var (
	router = mux.NewRouter()
)

func StartApplication() {
	elasticsearch.Init()

	mapUrls()

	srv := &http.Server{
		Addr:         "127.0.0.1:8082",
		Handler:      router,
		WriteTimeout: 500 * time.Millisecond,
		ReadTimeout:  3 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
