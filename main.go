package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	//	sellerDID []byte
}

//	type SellerReview struct {
//		DID          string  `cbor: "did"`
//		SellerRating float32 `cbor: "seller_rating"`
//		ProductCount float32 `cbor:"product_count"` //This can be changed to the list of products the seller has and calcuate the length of the array to get the no of products
//	}
type SellerReview struct {
	DID          string  `cbor:"did"`
	SellerRating float32 `cbor:"seller_rating"`
	ProductCount float32 `cbor:"product_count"`
}

// Suggestion : Each product can be made into a data token and the product description can be given as the content of the token
type ProductReview struct {
	ProductId   string  `cbor:"product_id"`
	Rating      float32 `cbor:"rating"`
	RatingCount float32 `cbor:"rating_count"`
	SellerDID   string  `cbor:"seller_did"`
}

func (r *WasmtimeRuntime) Init(wasmFile string) {
	fmt.Println("Initialising")
	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)
	err := linker.DefineWasi()
	fmt.Println("DefineWasi :", err)
	wasiConfig := wasmtime.NewWasiConfig()
	fmt.Println(wasiConfig)
	r.store = wasmtime.NewStore(engine)
	r.store.SetWasi(wasiConfig)
	linker.FuncWrap("env", "load_input", r.loadInput)
	linker.FuncWrap("env", "dump_output", r.dumpOutput)
	//	linker.FuncWrap("env", "get_account_info", r.getAccountInfo)
	wasmBytes, _ := os.ReadFile(wasmFile)
	module, _ := wasmtime.NewModule(r.store.Engine, wasmBytes)
	instance, _ := linker.Instantiate(r.store, module)
	r.memory = instance.GetExport(r.store, "memory").Memory()
	r.handler = instance.GetFunc(r.store, "handler")
}

func (r *WasmtimeRuntime) loadInput(pointer int32) {
	copy(r.memory.UnsafeData(r.store)[pointer:pointer+int32(len(r.input))], r.input)
}

/* Here when 2 strings are passed productId :"AB" and sellerDID : "DFG" the ouput obtained r,productId is "AB" and r.sellerdid is "ABD",
so the inference is when value is copied from r.store, irrespective of the pointer it is copying from the start of the memory till the length specified */

// Assuming `ProductReview` and `SellerReview` structs are defined correctly
func (r *WasmtimeRuntime) getAccountInfo(pointer int32, productStateLength int32, sellerStateLength int32) error {
	port := "20001"
	r.output = make([]byte, productStateLength+sellerStateLength)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+productStateLength+sellerStateLength])

	sellerReviewCbor := r.output[productStateLength:]
	fmt.Println("Seller Review CBOR :", sellerReviewCbor)
	sellerReview := SellerReview{}
	err := cbor.Unmarshal(sellerReviewCbor, &sellerReview)
	if err != nil {
		fmt.Println("Error unmarshaling SellerReview:", err)
	}
	fmt.Println("Seller DID :", sellerReview.DID)
	did := sellerReview.DID
	baseURL := fmt.Sprintf("http://localhost:%s/api/get-account-info", port)
	apiURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("Error parsing URL: %s", err)
	}

	// Add the query parameter to the URL
	queryValues := apiURL.Query()
	queryValues.Add("did", did)
	apiURL.RawQuery = queryValues.Encode()
	response, err := http.Get(apiURL.String())
	if err != nil {
		return fmt.Errorf("Error making GET request: %s", err)
	}
	fmt.Println("Response :", response)
	defer response.Body.Close()
	return nil
}

func (r *WasmtimeRuntime) dumpOutput(pointer int32, productReviewLength int32, sellerReviewLength int32) {
	fmt.Println("pointer:", pointer)
	fmt.Println("productReviewLength:", productReviewLength)
	fmt.Println("sellerReviewLength:", sellerReviewLength)

	r.output = make([]byte, productReviewLength+sellerReviewLength)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+productReviewLength+sellerReviewLength])

	sellerReviewCbor := r.output[productReviewLength:]
	fmt.Println("Seller Review CBOR :", sellerReviewCbor)
	sellerReview := SellerReview{}
	err := cbor.Unmarshal(sellerReviewCbor, &sellerReview)
	if err != nil {
		fmt.Println("Error unmarshaling SellerReview:", err)
	}
	fmt.Println("Seller DID :", sellerReview.DID)
	did := sellerReview.DID
	port := "20001"
	baseURL := fmt.Sprintf("http://localhost:%s/api/get-account-info", port)
	apiURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing URL: %s", err)
	}

	// Add the query parameter to the URL
	queryValues := apiURL.Query()
	queryValues.Add("did", did)
	queryValues.Add("port", port)

	apiURL.RawQuery = queryValues.Encode()
	response, err := http.Get(apiURL.String())
	if err != nil {
		fmt.Println("Error making GET request: %s", err)
	}
	fmt.Println("Response :", response)
	defer response.Body.Close()

	// Print the raw CBOR data to verify it is correct
	fmt.Println("Raw CBOR Data:", r.output)

	review := ProductReview{}
	//sellerReview := SellerReview{}
	cborData := r.output[:productReviewLength]
	cborDataSeller := r.output[productReviewLength:]
	fmt.Println("Length of cborData :", len(cborData))
	fmt.Println("Length of cborDataSeller :", len(cborDataSeller))
	// Print the CBOR data slices to verify they are correct
	fmt.Println("CBOR Data:", cborData)
	fmt.Printf("CBOR Data Seller: %v", cborDataSeller)

	err3 := cbor.Unmarshal(cborData, &review)
	if err3 != nil {
		fmt.Println("Error unmarshaling ProductReview:", err)
	}
	fmt.Println(review.ProductId)
	fmt.Println(review.Rating)
	fmt.Println(review.RatingCount)

	err2 := cbor.Unmarshal(cborDataSeller, &sellerReview)
	if err2 != nil {
		fmt.Println("Error unmarshaling SellerReview:", err2)
	}

	// Print the decoded values to verify they are correct
	fmt.Println("Latest Product Review:", review)
	fmt.Println("Latest Seller Review:", sellerReview)

	// Rest of the code for writing to JSON files
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

// TO DO : Optimise the memory usage
func (r *WasmtimeRuntime) RunHandler(data []byte, productStateLength int32, sellerStateLength int32, rating float32) []byte {
	r.input = data
	_, err := r.handler.Call(r.store, productStateLength, sellerStateLength, rating)
	if err != nil {
		panic(fmt.Errorf("Failed to call function: %v", err))
	}

	fmt.Println("Result:", r.output)
	return r.output
}

func ReadProductReview(filePath string) ProductReview {
	productStateUpdate, _ := os.ReadFile(filePath)
	var review ProductReview
	json.Unmarshal(productStateUpdate, &review)
	fmt.Println("ProductId : ", review.ProductId)
	fmt.Println("Current Rating : ", review.Rating)
	fmt.Println("Current Rating Count : ", review.RatingCount)
	fmt.Println("Seller DID : ", review.SellerDID)
	return review
}

func ReadSellerReview(filePath string) SellerReview {
	sellerStateUpdate, _ := os.ReadFile(filePath)
	var sellerReview SellerReview
	json.Unmarshal(sellerStateUpdate, &sellerReview)
	fmt.Println("SellerId: ", sellerReview.DID)
	fmt.Println("Seller Rating : ", sellerReview.SellerRating)
	fmt.Println("Product Count : ", sellerReview.ProductCount)
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
		panic(fmt.Errorf("Failed to encode string as CBOR: %v", err))
	}
	fmt.Println("Encoded Product State length :", len(encodedProductState))
	fmt.Println("CBOR encoded data :", encodedProductState)
	fmt.Println("ProductId : ", productStateUpdate.Rating)
	fmt.Println("Current Rating : ", productStateUpdate.Rating)
	fmt.Println("Current Rating Count : ", productStateUpdate.RatingCount)
	fmt.Println("Product Seller : ", productStateUpdate.SellerDID)
	//	randomRating := rand.Intn(5) + 1 //A random rating from 1-5 given for testing[Here it is considered as the rating a user gave]
	//whenever a new seller or product is registered
	randomRating := 5.00

	// sellerDID := "DFG"
	// sellerRating := float32(5)
	// productCount := float32(1)
	sellerStateUpdate := ReadSellerReview("store_state/rating_contract/seller_rating.json")

	fmt.Println("Random Rating :", randomRating)
	fmt.Println("SellerId: ", sellerStateUpdate.DID)
	fmt.Println("Seller Rating : ", sellerStateUpdate.SellerRating)
	fmt.Println("Product Count : ", sellerStateUpdate.ProductCount)

	encodedSellerState, err := cbor.Marshal(sellerStateUpdate)
	if err != nil {
		panic(fmt.Errorf("Failed to encode string as CBOR: %v", err))
	}
	fmt.Println("Encoded Seller State length :", len(encodedSellerState))
	review := ProductReview{}
	err3 := cbor.Unmarshal(encodedProductState, &review)
	if err3 != nil {
		fmt.Println("Error unmarshaling ProductReview:", err3)
	}

	fmt.Printf("%+v", review)

	fmt.Println("CBOR encoded data :", encodedSellerState)

	merge := append(encodedProductState, encodedSellerState...)
	fmt.Println("Merge : ", merge)
	fmt.Println("Merge length : ", len(merge))
	runtime := &WasmtimeRuntime{}
	runtime.Init("rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm")
	runtime.RunHandler(merge, int32(len(encodedProductState)), int32(len(encodedSellerState)), float32(randomRating))

}
