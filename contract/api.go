package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ContractInputRequest struct {
	SmartContractHash string `json:"smart_contract_hash"`
}

type ContractInputResponse struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

// mux route handler
func contractInputHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contract Input Handler")

	var req ContractInputRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	//should write a condition which checks whether the smart contract is already fetched or not , if fetched execute the contract directly,
	//or else if it is not fetched
	//fetch smart contract using fetchsmaetcontract api
	//Then the smart contract must be load, initialised and run

	resp := ContractInputResponse{Message: ""}
	json.NewEncoder(w).Encode(resp)

}
