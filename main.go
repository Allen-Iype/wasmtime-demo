package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/fxamacker/cbor/v2"
)

type WasmtimeRuntime struct {
	store   *wasmtime.Store
	memory  *wasmtime.Memory
	handler *wasmtime.Func

	input     []byte
	productId []byte
	sellerDID []byte
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
	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)
	linker.DefineWasi()
	wasiConfig := wasmtime.NewWasiConfig()
	r.store = wasmtime.NewStore(engine)
	r.store.SetWasi(wasiConfig)
	linker.FuncWrap("env", "load_input", r.loadInput)
	linker.FuncWrap("env", "dump_output", r.dumpOutput)
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

func (r *WasmtimeRuntime) dumpOutput(pointer int32, productReviewLength int32, sellerReviewLength int32) {
	fmt.Println("pointer:", pointer)
	fmt.Println("productReviewLength:", productReviewLength)
	fmt.Println("sellerReviewLength:", sellerReviewLength)

	r.productId = make([]byte, productReviewLength+sellerReviewLength)
	copy(r.productId, r.memory.UnsafeData(r.store)[pointer:pointer+productReviewLength+sellerReviewLength])

	// Print the raw CBOR data to verify it is correct
	fmt.Println("Raw CBOR Data:", r.productId)

	review := ProductReview{}
	sellerReview := SellerReview{}
	cborData := r.productId[:productReviewLength]
	cborDataSeller := r.productId[productReviewLength:]
	fmt.Println("Length of cborData :", len(cborData))
	fmt.Println("Length of cborDataSeller :", len(cborDataSeller))
	// Print the CBOR data slices to verify they are correct
	fmt.Println("CBOR Data:", cborData)
	fmt.Printf("CBOR Data Seller: %v", cborDataSeller)

	err := cbor.Unmarshal(cborData, &review)
	if err != nil {
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

// func (r *WasmtimeRuntime) dumpOutput(pointer int32, productReviewLength int32, sellerReviewLength int32) {
// 	fmt.Println("pointer :", pointer)
// 	fmt.Println("productReviewLength :", productReviewLength)
// 	fmt.Println("sellerReviewLength :", sellerReviewLength)

// 	r.productId = make([]byte, productReviewLength+sellerReviewLength)
// 	//r.sellerDID = make([]byte, sellerDidLength)
// 	copy(r.productId, r.memory.UnsafeData(r.store)[pointer:pointer+productReviewLength+sellerReviewLength])
// 	//	copy(r.sellerDID, r.memory.UnsafeData(r.store)[pointer:sellerDidPointer+sellerDidLength])
// 	//split byte array according to length
// 	review := ProductReview{}
// 	sellerReview := SellerReview{}
// 	cborData := r.productId[:productReviewLength]
// 	cborDataSeller := r.productId[productReviewLength:]
// 	fmt.Println("CborDta ", cborData)
// 	fmt.Println("CborDtaSeller ", cborDataSeller)

// 	err := cbor.Unmarshal(cborData, &review)
// 	fmt.Println("err :", err)
// 	if err != nil {
// 		// Handle the error appropriately, e.g., log the error, return an error response, etc.
// 		fmt.Println("Error unmarshaling ProductReview:", err)
// 	}
// 	err2 := cbor.Unmarshal(cborDataSeller, &sellerReview)
// 	if err2 != nil {
// 		// Handle the error appropriately, e.g., log the error, return an error response, etc.
// 		fmt.Println("Error unmarshaling SellerReview:", err)
// 	}
// 	fmt.Println("Latest Product Review", review)
// 	fmt.Println("Latest Seller Review", sellerReview)

// 	fmt.Println("Combined Byte Array :", r.productId)
// 	fmt.Println("Lenght of Byte Array", len(r.productId))
// 	content, err := json.Marshal(review)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = os.WriteFile("store_state/rating_contract/rating.json", content, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	sellerContent, err := json.Marshal(sellerReview)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = os.WriteFile("store_state/rating_contract/seller_rating.json", sellerContent, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// TO DO : Optimise the memory usage
func (r *WasmtimeRuntime) RunHandler(data []byte, productStateLength int32, sellerStateLength int32, rating float32) []byte {
	r.input = data
	_, err := r.handler.Call(r.store, productStateLength, sellerStateLength, rating)
	if err != nil {
		panic(fmt.Errorf("Failed to call function: %v", err))
	}

	fmt.Println("Result:", r.productId)
	return r.productId
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
