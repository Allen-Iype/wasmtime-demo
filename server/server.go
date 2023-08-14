package server

import (
	"fmt"
	"net/http"
	"wasmtime-demo/contract"

	"github.com/gorilla/mux"
)

func Bootup() {
	fmt.Println("Server Started")
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/contract-input", contract.ContractInputHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}
