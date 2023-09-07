package contract

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/joho/godotenv"
)

type WasmtimeRuntime struct {
	store   *wasmtime.Store
	memory  *wasmtime.Memory
	handler *wasmtime.Func

	input  []byte
	output []byte
}

type Count struct {
	Red           int32
	Blue          int32
	LatestBlockNo int32
}

type SmartContractDataReply struct {
	BasicResponse
	SCTDataReply []SCTDataReply
}

type BasicResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type SCTDataReply struct {
	BlockNo           uint64
	BlockId           string
	SmartContractData string
}

func (r *WasmtimeRuntime) loadInput(pointer int32) {
	copy(r.memory.UnsafeData(r.store)[pointer:pointer+int32(len(r.input))], r.input)
}

func (r *WasmtimeRuntime) Init(wasmFile string) {
	fmt.Println(wasmFile)
	fmt.Println("Initializing wasm")
	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)
	linker.DefineWasi()
	wasiConfig := wasmtime.NewWasiConfig()
	r.store = wasmtime.NewStore(engine)
	r.store.SetWasi(wasiConfig)
	linker.FuncWrap("env", "load_input", r.loadInput)
	linker.FuncWrap("env", "dump_output", r.dumpOutput)
	// linker.FuncWrap("env", "get_account_info", r.getAccountInfo)
	// linker.FuncWrap("env", "initiate_transfer", r.InitiateTransaction)
	wasmBytes, err := os.ReadFile(wasmFile)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %v", err))
	}
	module, _ := wasmtime.NewModule(r.store.Engine, wasmBytes)
	instance, _ := linker.Instantiate(r.store, module)
	r.memory = instance.GetExport(r.store, "memory").Memory()
	r.handler = instance.GetFunc(r.store, "handler")
}

// func (r *WasmtimeRuntime) getAccountInfo() {
// 	fmt.Println("Get Account Info")
// 	port := "20002"
// 	productReviewLength := 59 //issue here
// 	sellerReviewCbor := r.output[productReviewLength:]
// 	fmt.Println("Seller Review CBOR :", sellerReviewCbor)
// 	sellerReview := SellerReview{}
// 	err := cbor.Unmarshal(sellerReviewCbor, &sellerReview)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling SellerReview:", err)
// 	}
// 	fmt.Println("Seller DID :", sellerReview.DID)
// 	did := sellerReview.DID
// 	//	did := "bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
// 	baseURL := fmt.Sprintf("http://localhost:%s/api/get-account-info", port)
// 	apiURL, err := url.Parse(baseURL)
// 	fmt.Println(apiURL)
// 	if err != nil {
// 		fmt.Printf("Error parsing URL: %s\n", err)
// 		return
// 	}

// 	// Add the query parameter to the URL
// 	queryValues := apiURL.Query()
// 	queryValues.Add("did", did)
// 	queryValues.Add("port", port)
// 	fmt.Println("Query Values", queryValues)
// 	apiURL.RawQuery = queryValues.Encode()
// 	fmt.Println("Api Raw Query URL:", apiURL.RawQuery)
// 	fmt.Println("Query Values Encode:", queryValues.Encode())
// 	fmt.Println("Api URL string:", apiURL.String())
// 	response, err := http.Get(apiURL.String())
// 	if err != nil {
// 		fmt.Printf("Error making GET request: %s\n", err)
// 		return
// 	}
// 	fmt.Println("Response Status:", response.Status)
// 	defer response.Body.Close()

// 	// Handle the response data as needed
// 	if response.StatusCode == http.StatusOK {
// 		data, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			fmt.Printf("Error reading response body: %s\n", err)
// 			return
// 		}
// 		// Process the data as needed
// 		fmt.Println("Response Body:", string(data))
// 	} else {
// 		fmt.Printf("API returned a non-200 status code: %d\n", response.StatusCode)
// 		data, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			fmt.Printf("Error reading error response body: %s\n", err)
// 			return
// 		}
// 		fmt.Println("Error Response Body:", string(data))
// 		return
// 	}
// }

func (r *WasmtimeRuntime) InitiateTransaction() {
	port := "20002"
	receiver := "12D3KooWSokjA3JcWZNJUz4B6mN7tYBH75bSSGoxQwqJ1kTBSvgM.bafybmiegyiz5zvnveqx3lc3cealx3zfwiclwpntaf3ep3zm2slexbzj33u"
	sender := "12D3KooWCR4BW7gfPmCZhAJusqv1PoS49jgqTGvofcG4WPyg8FxV.bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
	tokenCount := 1
	comment := "Wasm Test"

	data := map[string]interface{}{
		"receiver":   receiver,
		"sender":     sender,
		"tokenCOunt": tokenCount,
		"comment":    comment,
		"type":       2,
	}

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println("initiateTransactionPayload request to rubix:", string(bodyJSON))

	url := fmt.Sprintf("http://localhost:%s/api/initiate-rbt-transfer", port)
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
	fmt.Println("Response Body:", string(data2))

	defer resp.Body.Close()
}

func (r *WasmtimeRuntime) RunHandler(data []byte, inputVoteLength int32, redLength int32, blueLength int32, portLength int32, hashLength int32, blockNoLength int32) []byte {
	r.input = data
	_, err := r.handler.Call(r.store, inputVoteLength, redLength, blueLength, portLength, hashLength, blockNoLength)
	if err != nil {
		panic(fmt.Errorf("failed to call function: %v", err))
	}
	return r.output
}

func (r *WasmtimeRuntime) dumpOutput(pointer int32, red int32, blue int32, block_no int32, port_length int32, hash_length int32) {
	fmt.Println("red :", red)
	fmt.Println("blue :", blue)
	fmt.Println("port_length :", port_length)
	fmt.Println("hash_length :", hash_length)
	fmt.Println("block_no :", block_no)
	fmt.Println("pointer :", pointer)
	r.output = make([]byte, port_length+hash_length)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+(port_length+hash_length)])
	fmt.Println("output array :", r.output)
	err3 := godotenv.Load()
	if err3 != nil {
		fmt.Println("Error loading .env file:", err3)
		return
	}
	port := string(r.output[:port_length])
	smartContracthash := string(r.output[port_length:])
	// port := string(r.output[:port_length])
	fmt.Println("The port is :", port)
	// smartContractHash := string(r.output[port_length : p])
	fmt.Println("Smart Contract Hash in Dump Output :", smartContracthash)
	nodeName := os.Getenv(port)
	stateFilePath := fmt.Sprintf("/home/allen/Rubix-Wasm_test/%s/SmartContract/%s/schemaCodeFile.json", nodeName, smartContracthash)

	count := Count{}
	count.Red = int32(red)
	count.Blue = int32(blue)
	count.LatestBlockNo = int32(block_no)

	fmt.Println("Count in Dump Output :", count)

	content, err := json.Marshal(count)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Json Marshalled Count :", content)
	err = os.WriteFile(stateFilePath, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateSmartContract(did string, wasmPath string, schemaPath string, rawCodePath string, port string) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the form fields
	_ = writer.WriteField("did", did)

	// Add the binaryCodePath field
	file, _ := os.Open(wasmPath)
	defer file.Close()
	binaryPart, _ := writer.CreateFormFile("binaryCodePath", wasmPath)
	_, _ = io.Copy(binaryPart, file)

	// Add the rawCodePath field
	rawFile, _ := os.Open(rawCodePath)
	defer rawFile.Close()
	rawPart, _ := writer.CreateFormFile("rawCodePath", rawCodePath)
	_, _ = io.Copy(rawPart, rawFile)

	// Add the schemaFilePath field
	schemaFile, _ := os.Open(schemaPath)
	defer schemaFile.Close()
	schemaPart, _ := writer.CreateFormFile("schemaFilePath", schemaPath)
	_, _ = io.Copy(schemaPart, schemaFile)

	// Close the writer
	writer.Close()

	// Create the HTTP request
	url := fmt.Sprintf("http://localhost:%s/api/generate-smart-contract", port)
	req, _ := http.NewRequest("POST", url, &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in execute Contract :", string(data2))

	// Process the response as needed
	fmt.Println("Response status code:", resp.StatusCode)
}

func GetSmartContractData(port string, token string) []byte {
	data := map[string]interface{}{
		"token":  token,
		"latest": false,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}
	url := fmt.Sprintf("http://localhost:%s/api/get-smart-contract-token-chain-data", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
	}
	// Process the data as needed
	fmt.Println("Response Body in get smart contract data :", string(data2))

	return data2

}

func DeploySmartContract(comment string, deployerAddress string, quorumType int, rbtAmount int, smartContractToken string, port string) string {
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
	}
	url := fmt.Sprintf("http://localhost:%s/api/deploy-smart-contract", port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
	}
	// Process the data as needed
	fmt.Println("Response Body in deploy smart contract:", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)

	defer resp.Body.Close()
	return id

}

func SignatureResponse(requestId string, port string) {
	data := map[string]interface{}{
		"id":       requestId,
		"mode":     0,
		"password": "mypassword",
	}

	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/signature-response", port)
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
	fmt.Println("Response Body in signature response :", string(data2))
	//json encode string
	defer resp.Body.Close()

}

func ExecuteSmartContract(comment string, executorAddress string, quorumType int, smartContractData string, smartContractToken string, port string) {
	data := map[string]interface{}{
		"comment":            comment,
		"executorAddr":       executorAddress,
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
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	SignatureResponse(id, port)

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
	url := fmt.Sprintf("http://localhost:%s/api/subscribe-contract", port)
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
	fmt.Println("Response Body in subscribe smart contract :", string(data2))

	defer resp.Body.Close()

}

func FetchSmartContract(smartContractTokenHash string, port string) {
	data := map[string]interface{}{
		"smart_contract_token": smartContractTokenHash,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/fetch-smart-contract", port)
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
	fmt.Println("Response Body in fetch smart contract :", string(data2))

	defer resp.Body.Close()

}

func RegisterCallBackUrl(smartContractTokenHash string, urlPort string, endPoint string, nodePort string) {
	callBackUrl := fmt.Sprintf("http://localhost:%s/%s", urlPort, endPoint)
	data := map[string]interface{}{
		"CallBackURL":        callBackUrl,
		"SmartContractToken": smartContractTokenHash,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/register-callback-url", nodePort)
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
	fmt.Println("Response Body in register callback url :", string(data2))
}

// func ReadProductReview(filePath string) ProductReview {
// 	productStateUpdate, _ := os.ReadFile(filePath)
// 	var review ProductReview
// 	json.Unmarshal(productStateUpdate, &review)
// 	return review
// }

// func ReadSellerReview(filePath string) SellerReview {
// 	sellerStateUpdate, _ := os.ReadFile(filePath)
// 	var sellerReview SellerReview
// 	json.Unmarshal(sellerStateUpdate, &sellerReview)
// 	return sellerReview
// }

func ReadCurrentState(stateFilePath string) string {
	currentStateJsonFile, err := os.ReadFile(stateFilePath)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	currentState := string(currentStateJsonFile)
	return currentState
}

func GetRubixSmartContractPath(contractHash string, smartContractName string, nodeName string) (string, error) {
	rubixcontractPath := fmt.Sprintf("/home/allen/Rubix-Wasm_test/%s/SmartContract/%s/%s", nodeName, contractHash, smartContractName)

	// Check if the path exists
	if _, err := os.Stat(rubixcontractPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Smart contract path does not exist")
		}
		return "", err // Return other errors as is
	}

	return rubixcontractPath, nil
}

func GetRubixSchemaPath(contractHash string, nodeName string, schemaName string) (string, error) {
	rubixSchemaPath := fmt.Sprintf("/home/allen/Rubix-Wasm_test/%s/SmartContract/%s/%s", nodeName, contractHash, schemaName)

	// Check if the path exists
	if _, err := os.Stat(rubixSchemaPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Schema path does not exist")
		}
		return "", err // Return other errors as is
	}

	return rubixSchemaPath, nil
}

func WasmInput() {

}

func RunSmartContract(wasmPath string, schemaPath string, port string, smartContractTokenHash string) {

	smartContractTokenData := GetSmartContractData(port, smartContractTokenHash)
	fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))

	var dataReply SmartContractDataReply

	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Data reply in RunSmartContract", dataReply)
	runtime := &WasmtimeRuntime{}
	//runtime.Init("rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm")
	runtime.Init(wasmPath)
	//While this loop is running, there is a question whether any state condition needs to be checked at this point.

	// instead of the runhandler calling all the inputs, a save state function must be created just to update the state, rest everything should
	//be handled by the runhandler

	// Process each SCTDataReply item in the array
	// Check BlockId from the end of the array

	var count Count

	byteValue, _ := os.ReadFile(schemaPath)
	json.Unmarshal(byteValue, &count)

	//	targetBlockNo := int(count.LatestBlockNo)
	targetBlockNo := count.LatestBlockNo

	//	previousBlock := count.LatestBlockHash
	redvote := count.Red
	bluevote := count.Blue

	//sctReply := dataReply.SCTDataReply
	// Check if the target block number is within the bounds of the array
	if int(targetBlockNo) < len(dataReply.SCTDataReply) {
		// Perform operations on elements after the target block
		for i := int(targetBlockNo) + 1; i < len(dataReply.SCTDataReply); i++ {
			// Access and operate on reply.SCTDataReply[i]
			fmt.Println("Performing operation on BlockNo:", dataReply.SCTDataReply[i].BlockNo)

			// if dataReply.SCTDataReply[i].BlockNo == 0 {
			// 	continue // Skip this iteration and proceed to the next one when the block number is zero
			// }

			//	if dataReply.SCTDataReply[i].BlockNo == 0 || dataReply.SCTDataReply[i-1].BlockId == previousBlock {
			fmt.Println("previous block is same")

			inputVote := []byte(dataReply.SCTDataReply[i].SmartContractData)
			fmt.Println("Input vote is :", inputVote)
			inputBlockId := []byte(dataReply.SCTDataReply[i].BlockId)
			fmt.Println("Input block id is :", inputBlockId)
			inputBlockNo := make([]byte, 4)
			binary.LittleEndian.PutUint32(inputBlockNo, uint32(dataReply.SCTDataReply[i].BlockNo))

			red := make([]byte, 4)
			binary.LittleEndian.PutUint32(red, uint32(redvote))

			blue := make([]byte, 4)
			binary.LittleEndian.PutUint32(blue, uint32(bluevote))

			portByte := []byte(port)

			smartContractHashByte := []byte(smartContractTokenHash)

			portAndHash := append(portByte, smartContractHashByte...)

			mergevote := append(red, blue...)
			fmt.Println("mergevote ", mergevote)

			merge := append(inputVote, mergevote...)

			mergePortAndHash := append(merge, portAndHash...)

			//	blockIdAndNo := append(inputBlockId, inputBlockNo...)

			mergeComplete := append(mergePortAndHash, inputBlockNo...)

			fmt.Println("merge complete", mergeComplete)

			runtime.RunHandler(mergeComplete, int32(len(inputVote)), int32(len(red)), int32(len(blue)), int32(len(portByte)), int32(len(smartContractHashByte)), int32(len(inputBlockNo)))

			// } else {
			// 	fmt.Println("previous block is not matching")

			// }

			// Perform your operations here
		}
	} else {
		fmt.Println("Target block number is out of bounds")
	}

}

// func RunSmartContract(wasmPath string, port string, smartContractTokenHash string) {

// 	smartContractTokenData := GetSmartContractData(port, smartContractTokenHash)
// 	fmt.Println("Smart Contract Token Data :", string(smartContractTokenData))

// 	var dataReply SmartContractDataReply

// 	if err := json.Unmarshal(smartContractTokenData, &dataReply); err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Data reply in RunSmartContract", dataReply)
// 	runtime := &WasmtimeRuntime{}
// 	//runtime.Init("rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm")
// 	runtime.Init(wasmPath)
// 	//While this loop is running, there is a question whether any state condition needs to be checked at this point.

// 	// instead of the runhandler calling all the inputs, a save state function must be created just to update the state, rest everything should
// 	//be handled by the runhandler

// 	// Process each SCTDataReply item in the array
// 	for _, sctReply := range dataReply.SCTDataReply {
// 		fmt.Println("BlockNo:", sctReply.BlockNo)
// 		fmt.Println("BlockId:", sctReply.BlockId)
// 		fmt.Println("SmartContractData:", sctReply.SmartContractData)
// 		productStateUpdate := ReadProductReview("store_state/rating_contract/rating.json")
// 		encodedProductState, err := cbor.Marshal(productStateUpdate)
// 		if err != nil {
// 			panic(fmt.Errorf("failed to encode string as CBOR: %v", err))
// 		}
// 		//	randomRating := rand.Intn(5) + 1 //A random rating from 1-5 given for testing[Here it is considered as the rating a user gave]
// 		//whenever a new seller or product is registered
// 		randomRating := sctReply.SmartContractData
// 		floatValue, err := strconv.ParseFloat(randomRating, 32)
// 		if err != nil {
// 			fmt.Println("Error converting string to float:", err)
// 		}
// 		sellerStateUpdate := ReadSellerReview("store_state/rating_contract/seller_rating.json")
// 		fmt.Println("Random Rating :", randomRating)
// 		fmt.Println("SellerId: ", sellerStateUpdate.DID)
// 		fmt.Println("Seller Rating : ", sellerStateUpdate.SellerRating)
// 		fmt.Println("Product Count : ", sellerStateUpdate.ProductCount)

// 		encodedSellerState, err := cbor.Marshal(sellerStateUpdate)
// 		if err != nil {
// 			panic(fmt.Errorf("failed to encode string as CBOR: %v", err))
// 		}
// 		review := ProductReview{}
// 		err3 := cbor.Unmarshal(encodedProductState, &review)
// 		if err3 != nil {
// 			fmt.Println("Error unmarshaling ProductReview:", err3)
// 		}

// 		fmt.Printf("%+v", review)

// 		fmt.Println("CBOR encoded data :", encodedSellerState)

// 		merge := append(encodedProductState, encodedSellerState...)
// 		runtime.RunHandler(merge, int32(len(encodedProductState)), int32(len(encodedSellerState)), float32(floatValue))
// 		// Perform your operations on each sctReply item here
// 		fmt.Println()
// 	}

// }
