package contract

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ContractInputRequest struct {
	Port              string `json:"port"`
	SmartContractHash string `json:"smart_contract_hash"` //port should also be added here, so that the api can understand which node.
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
	err3 := godotenv.Load()
	if err3 != nil {
		fmt.Println("Error loading .env file:", err3)
		return
	}
	port := req.Port
	nodeName := os.Getenv(port)
	folderPath, _ := GetRubixSmartContractPath(req.SmartContractHash, "binaryCodeFile.wasm", nodeName)
	schemaPath, _ := GetRubixSchemaPath(req.SmartContractHash, nodeName, "schemaCodeFile.json")
	fmt.Println(folderPath)
	_, err1 := os.Stat(folderPath)
	fmt.Println(err1)
	if os.IsNotExist(err1) {
		fmt.Println("Smart Contract not found")
		//FetchSmartContract(req.SmartContractHash, port)
		RunSmartContract(folderPath, schemaPath, port, req.SmartContractHash)
	} else if err == nil {
		fmt.Printf("Folder '%s' exists", folderPath)

		RunSmartContract(folderPath, schemaPath, port, req.SmartContractHash)

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
