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
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
