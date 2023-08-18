package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ContractInputRequest struct {
	SmartContractHash string `json:"smart_contract_hash"`
}

type ContractInputResponse struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

// /api/get-smart-contract-data
// latest: false
//
// mux route handler
func ContractInputHandler(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Contract Input Handler")

	var req ContractInputRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	//TODO: Smart contract should be fetched into the dapp folder or an api must be provided to get the smart contract path
	// As of now we are hardcoding the port as well as the smart contract path is fetched using the GetRubixSmartContractPath function
	port := "20002"
	folderPath := GetRubixSmartContractPath(req.SmartContractHash)
	_, err1 := os.Stat(folderPath)
	if os.IsNotExist(err1) {
		fmt.Println("Smart Contract not found")
		FetchSmartContract(req.SmartContractHash, port)
		RunSmartContract(folderPath, port)
	} else if err == nil {
		fmt.Printf("Folder '%s' exists", folderPath)

		RunSmartContract(folderPath, port)

	} else {
		fmt.Printf("Error while checking folder: %v\n", err)
	}

	//should write a condition which checks whether the smart contract is already fetched or not , if fetched execute the contract directly,
	//or else if it is not fetched
	//fetch smart contract using fetchsmaetcontract api
	//Then the smart contract must be load, initialised and run

	resp := ContractInputResponse{Message: "Callback Successful", Result: "Success"}
	json.NewEncoder(w).Encode(resp)

}
