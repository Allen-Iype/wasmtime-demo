package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// this hsould include the init and reading of the smart contract part
func GenerateSmartContract(did string, wasmPath string, schemaPath string, rawCodePath string, port string) {
	url := fmt.Sprintf("http://localhost:%s/api/generate-smart-contract", port)

	// Create a new buffer to write the multipart request
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the fields to the request
	multipartWriter.WriteField("did", did)
	multipartWriter.CreateFormFile("binaryCodePath", wasmPath)
	multipartWriter.CreateFormFile("rawCodePath", rawCodePath)
	multipartWriter.CreateFormFile("schemaFilePath", schemaPath)

	// Close the multipart writer to finalize the request
	multipartWriter.Close()

	// Create the HTTP request
	request, _ := http.NewRequest("POST", url, &requestBody)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}
	defer response.Body.Close()
	data2, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute Contract :", string(data2))

	// Process the response as needed
	fmt.Println("Response status code:", response.StatusCode)

}

func DeploySmartContract(comment string, deployerAddress string, quorumType int, rbtAmount int, smartContractToken string, port string) {
	data := map[string]interface{}{
		"comment":            comment,
		"deployerAddr":       deployerAddress,
		"quorumType":         quorumType,
		"rbtAmount":          rbtAmount,
		"smartContractToken": smartContractToken,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/deploy-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in deploy smart contract:", string(data2))

	defer resp.Body.Close()

}

func ExecuteSmartContract(comment string, executorAddress string, quorumType int, smartContractData string, smartContractToken string, port string) {
	data := map[string]interface{}{
		"comment":            comment,
		"executorAddress":    executorAddress,
		"quorumType":         quorumType,
		"smartContractData":  smartContractData,
		"smartContractToken": smartContractToken,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/execute-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute smart contract :", string(data2))

	defer resp.Body.Close()

}

func SubscribeSmartContract(contractToken string, port string) {
	data := map[string]interface{}{
		"contract": contractToken,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/subscribe-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute smart contract :", string(data2))

	defer resp.Body.Close()

}

func FetchSmartContract(smartContractTokenHash string, port string) {

}

func RunSmartContract() {

}
