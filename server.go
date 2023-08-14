package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func callBackapi() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/health", healthHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}
