package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/fxamacker/cbor/v2"
)

type WasmtimeRuntime struct {
	store   *wasmtime.Store
	memory  *wasmtime.Memory
	handler *wasmtime.Func

	input  []byte
	output []byte
}

/*From the discussions what I understood is since the Rust part or the Smart contract is written by the DAPP, they would know the structure
of the data which is being passed or the data which is being used, so this structure which the smart contract deployer or the dapp has
created must be known to the Go-package. So this structure which is being written in the Rust code must be provided to the GO-code
as the Schema  */

type SellerReview struct {
	DID          string  `cbor:"did"`
	SellerRating float32 `cbor:"seller_rating"`
	ProductCount float32 `cbor:"product_count"` //This can be changed to the list of products the seller has and calcuate the length of the array to get the no of products
}

// Suggestion : Each product can be made into a data token and the product description can be given as the content of the token
type ProductReview struct {
	ProductId   string  `cbor:"product_id"`
	Rating      float32 `cbor:"rating"`
	RatingCount float32 `cbor:"rating_count"`
	SellerDID   string  `cbor:"seller_did"`
}

func (r *WasmtimeRuntime) Init(wasmFile string) {
	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)
	linker.DefineWasi()
	wasiConfig := wasmtime.NewWasiConfig()
	r.store = wasmtime.NewStore(engine)
	r.store.SetWasi(wasiConfig)
	linker.FuncWrap("env", "load_input", r.loadInput)
	linker.FuncWrap("env", "dump_output", r.dumpOutput)
	linker.FuncWrap("env", "get_account_info", r.getAccountInfo)
	linker.FuncWrap("env", "initiate_transfer", r.InitiateTransaction)
	wasmBytes, _ := os.ReadFile(wasmFile)
	module, _ := wasmtime.NewModule(r.store.Engine, wasmBytes)
	instance, _ := linker.Instantiate(r.store, module)
	r.memory = instance.GetExport(r.store, "memory").Memory()
	r.handler = instance.GetFunc(r.store, "handler")
}

func (r *WasmtimeRuntime) loadInput(pointer int32) {
	copy(r.memory.UnsafeData(r.store)[pointer:pointer+int32(len(r.input))], r.input)
}

//schema should be created
//struct and things are defined in Rust, this schema must be provided to golang for identifying the structure
//basically the state is stored in the tokenchain.

func (r *WasmtimeRuntime) dumpOutput(pointer int32, productReviewLength int32, sellerReviewLength int32) {

	r.output = make([]byte, productReviewLength+sellerReviewLength)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+productReviewLength+sellerReviewLength])

	review := ProductReview{}
	sellerReview := SellerReview{}
	cborData := r.output[:productReviewLength]
	cborDataSeller := r.output[productReviewLength:]

	err3 := cbor.Unmarshal(cborData, &review)
	if err3 != nil {
		fmt.Println("Error unmarshaling ProductReview:", err3)
	}
	fmt.Println(review.ProductId)
	fmt.Println(review.Rating)
	fmt.Println(review.RatingCount)

	err2 := cbor.Unmarshal(cborDataSeller, &sellerReview)
	if err2 != nil {
		fmt.Println("Error unmarshaling SellerReview:", err2)
	}

	content, err := json.Marshal(review)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("store_state/rating_contract/rating.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	sellerContent, err := json.Marshal(sellerReview)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("store_state/rating_contract/seller_rating.json", sellerContent, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *WasmtimeRuntime) getAccountInfo() {
	fmt.Println("Get Account Info")
	port := "20002"
	productReviewLength := 59 //issue here
	sellerReviewCbor := r.output[productReviewLength:]
	fmt.Println("Seller Review CBOR :", sellerReviewCbor)
	sellerReview := SellerReview{}
	err := cbor.Unmarshal(sellerReviewCbor, &sellerReview)
	if err != nil {
		fmt.Println("Error unmarshaling SellerReview:", err)
	}
	fmt.Println("Seller DID :", sellerReview.DID)
	did := sellerReview.DID
	//	did := "bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
	baseURL := fmt.Sprintf("http://localhost:%s/api/get-account-info", port)
	apiURL, err := url.Parse(baseURL)
	fmt.Println(apiURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %s\n", err)
		return
	}

	// Add the query parameter to the URL
	queryValues := apiURL.Query()
	queryValues.Add("did", did)
	queryValues.Add("port", port)
	fmt.Println("Query Values", queryValues)
	apiURL.RawQuery = queryValues.Encode()
	fmt.Println("Api Raw Query URL:", apiURL.RawQuery)
	fmt.Println("Query Values Encode:", queryValues.Encode())
	fmt.Println("Api URL string:", apiURL.String())
	response, err := http.Get(apiURL.String())
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		return
	}
	fmt.Println("Response Status:", response.Status)
	defer response.Body.Close()

	// Handle the response data as needed
	if response.StatusCode == http.StatusOK {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err)
			return
		}
		// Process the data as needed
		fmt.Println("Response Body:", string(data))
	} else {
		fmt.Printf("API returned a non-200 status code: %d\n", response.StatusCode)
		data, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading error response body: %s\n", err)
			return
		}
		fmt.Println("Error Response Body:", string(data))
		return
	}
}

func (r *WasmtimeRuntime) InitiateTransaction(reveiverinput string) {
	port := "20002"
	receiver := "12D3KooWSokjA3JcWZNJUz4B6mN7tYBH75bSSGoxQwqJ1kTBSvgM.bafybmiegyiz5zvnveqx3lc3cealx3zfwiclwpntaf3ep3zm2slexbzj33u"
	sender := "12D3KooWCR4BW7gfPmCZhAJusqv1PoS49jgqTGvofcG4WPyg8FxV.bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
	tokenCount := 1
	comment := "Wasm Test"

	data := map[string]interface{}{
		reveiverinput: receiver,
		"sender":      sender,
		"tokenCOunt":  tokenCount,
		"comment":     comment,
		"type":        2,
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

// TO DO : Optimise the memory usage
func (r *WasmtimeRuntime) RunHandler(data []byte, productStateLength int32, sellerStateLength int32, rating float32) []byte {
	r.input = data
	_, err := r.handler.Call(r.store, productStateLength, sellerStateLength, rating)
	if err != nil {
		panic(fmt.Errorf("failed to call function: %v", err))
	}
	return r.output
}

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

func ReadProductReview(filePath string) ProductReview {
	productStateUpdate, _ := os.ReadFile(filePath)
	var review ProductReview
	json.Unmarshal(productStateUpdate, &review)
	return review
}

func ReadSellerReview(filePath string) SellerReview {
	sellerStateUpdate, _ := os.ReadFile(filePath)
	var sellerReview SellerReview
	json.Unmarshal(sellerStateUpdate, &sellerReview)
	return sellerReview
}

func ConvertFloat32ToBytes(floatValue float32) []byte {
	bits := math.Float32bits(floatValue)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func main() {
	productStateUpdate := ReadProductReview("store_state/rating_contract/rating.json")
	encodedProductState, err := cbor.Marshal(productStateUpdate)
	if err != nil {
		panic(fmt.Errorf("failed to encode string as CBOR: %v", err))
	}
	//	randomRating := rand.Intn(5) + 1 //A random rating from 1-5 given for testing[Here it is considered as the rating a user gave]
	//whenever a new seller or product is registered
	randomRating := 5.00

	sellerStateUpdate := ReadSellerReview("store_state/rating_contract/seller_rating.json")
	fmt.Println("Random Rating :", randomRating)
	fmt.Println("SellerId: ", sellerStateUpdate.DID)
	fmt.Println("Seller Rating : ", sellerStateUpdate.SellerRating)
	fmt.Println("Product Count : ", sellerStateUpdate.ProductCount)

	encodedSellerState, err := cbor.Marshal(sellerStateUpdate)
	if err != nil {
		panic(fmt.Errorf("failed to encode string as CBOR: %v", err))
	}
	review := ProductReview{}
	err3 := cbor.Unmarshal(encodedProductState, &review)
	if err3 != nil {
		fmt.Println("Error unmarshaling ProductReview:", err3)
	}

	fmt.Printf("%+v", review)

	fmt.Println("CBOR encoded data :", encodedSellerState)

	merge := append(encodedProductState, encodedSellerState...)
	runtime := &WasmtimeRuntime{}
	runtime.Init("rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm")
	runtime.RunHandler(merge, int32(len(encodedProductState)), int32(len(encodedSellerState)), float32(randomRating))

}
