package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
)

type WasmtimeRuntime struct {
	store   *wasmtime.Store
	memory  *wasmtime.Memory
	handler *wasmtime.Func

	input     []byte
	productId []byte
	sellerDID []byte
}

type SellerReview struct {
	DID          string
	SellerRating float32
	ProductCount float32 //This can be changed to the list of products the seller has and calcuate the length of the array to get the no of products
}

// Suggestion : Each product can be made into a data token and the product description can be given as the content of the token
type ProductReview struct {
	ProductId   string
	Rating      float32
	RatingCount float32
	SellerDID   string
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

func (r *WasmtimeRuntime) dumpOutput(pointer int32, latestRating float32, ratingCount float32, productIdLength int32, sellerDidLength int32, currentSellerRating float32, sellerDidPointer int32) {
	fmt.Println("latestRating :", latestRating)
	fmt.Println("ratingCount :", ratingCount)
	fmt.Println("productIdLength :", productIdLength)
	fmt.Println("sellerDidLength :", sellerDidLength)
	fmt.Println("currentSellerRating :", currentSellerRating)
	fmt.Println("sellerDidPointer :", sellerDidPointer)

	r.productId = make([]byte, productIdLength)
	r.sellerDID = make([]byte, sellerDidLength)
	fmt.Println(sellerDidLength)
	fmt.Println(productIdLength)
	copy(r.productId, r.memory.UnsafeData(r.store)[pointer:pointer+productIdLength])
	copy(r.sellerDID, r.memory.UnsafeData(r.store)[pointer:sellerDidPointer+sellerDidLength])
	review := ProductReview{}
	fmt.Println("Product id :", r.productId)
	fmt.Println("Lenght of r.productId", len(r.productId))
	review.ProductId = string(r.productId)
	review.Rating = float32(latestRating)
	review.RatingCount = float32(ratingCount)
	content, err := json.Marshal(review)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("store_state/rating_contract/rating.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
	sellerReview := SellerReview{}
	fmt.Println("Product Id pointer :", pointer)
	fmt.Println("Seller Did pointer :", sellerDidPointer)
	fmt.Println("Seller Did :", r.sellerDID)
	fmt.Println("Length of r.sellerdid", len(r.sellerDID))
	fmt.Println("Seller did string: ", string(r.sellerDID))
	fmt.Println("Product Id string: ", string(r.productId))
	sellerReview.DID = string(r.sellerDID)
	sellerReview.SellerRating = float32(currentSellerRating)
	sellerReview.ProductCount = float32(1)
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
func (r *WasmtimeRuntime) RunHandler(data []byte, didLength int32, ratingLength int32, countLength int32, userRatingLength int32, sellerDidLength int32, sellerRatingLength int32, sellerProductCountLength int32) []byte {
	r.input = data
	r.handler.Call(r.store, didLength, ratingLength, countLength, userRatingLength, sellerDidLength, sellerRatingLength, sellerProductCountLength)
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
	// productStateUpdate := ReadProductReview("store_state/rating_contract/rating.json")
	// currentRating := productStateUpdate.Rating
	// ratingCount := productStateUpdate.RatingCount
	// productId := productStateUpdate.ProductId

	productId := "AB"
	currentRating := float32(5)
	ratingCount := float32(1)
	productSeller := "DFG"

	fmt.Println("ProductId : ", productId)
	fmt.Println("Current Rating : ", currentRating)
	fmt.Println("Current Rating Count : ", ratingCount)
	fmt.Println("Product Seller : ", productSeller)
	//	randomRating := rand.Intn(5) + 1 //A random rating from 1-5 given for testing[Here it is considered as the rating a user gave]
	//whenever a new seller or product is registered
	randomRating := 5

	sellerDID := "DFG"
	sellerRating := float32(5)
	productCount := float32(1)
	// sellerStateUpdate := ReadSellerReview("store_state/rating_contract/seller_rating.json")
	// sellerRating := sellerStateUpdate.SellerRating
	// productCount := sellerStateUpdate.ProductCount
	// sellerDID := sellerStateUpdate.DID
	fmt.Println("Random Rating :", randomRating)
	fmt.Println("SellerId: ", sellerDID)
	fmt.Println("Seller Rating : ", sellerRating)
	fmt.Println("Product Count : ", productCount)

	productIdBytes := []byte(productId)
	fmt.Println("ProductIdBytes :", productIdBytes)
	fmt.Println("Length of ProductIdBytes :", len(productIdBytes))

	//the current average rating of the product
	ratingBytes := ConvertFloat32ToBytes(currentRating)

	//the total count of the ratings received
	countBytes := ConvertFloat32ToBytes(ratingCount)

	//the new rating given by the user
	userRatingBytes := ConvertFloat32ToBytes(float32(randomRating))

	sellerDIDBytes := []byte(sellerDID)
	fmt.Println("SellerDID bytes :", sellerDIDBytes)
	fmt.Println("Length of SellerDID bytes :", len(sellerDIDBytes))
	sellerRatingBytes := ConvertFloat32ToBytes(sellerRating)

	sellerProductCountBytes := ConvertFloat32ToBytes(productCount)

	productIdSellerDid := append(productIdBytes, sellerDIDBytes...)
	fmt.Println("Product Id and Seller Did :", productIdSellerDid)

	mergeSellerRating := append(productIdSellerDid, sellerRatingBytes...)
	fmt.Println("Product Id, Seller Did and Seller Rating :", mergeSellerRating)

	mergeSellerProductCount := append(mergeSellerRating, sellerProductCountBytes...)
	fmt.Println("Product Id, Seller Did, Seller Rating and Seller Product Count :", mergeSellerProductCount)

	mergeUserRating := append(mergeSellerProductCount, userRatingBytes...)
	fmt.Println("Product Id, Seller Did, Seller Rating, Seller Product Count and User Rating :", mergeUserRating)

	mergeCount := append(mergeUserRating, countBytes...)
	fmt.Println("Product Id, Seller Did, Seller Rating, Seller Product Count, User Rating and Count :", mergeCount)

	merge := append(mergeCount, ratingBytes...)
	fmt.Println("Product Id, Seller Did, Seller Rating, Seller Product Count, User Rating, Count and Rating :", merge)

	runtime := &WasmtimeRuntime{}
	runtime.Init("rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm")
	runtime.RunHandler(merge, int32(len(productIdBytes)), int32(len(ratingBytes)), int32(len(countBytes)), int32(len(userRatingBytes)), int32(len(sellerDIDBytes)), int32(len(sellerRatingBytes)), int32(len(sellerProductCountBytes)))

	// mergeCurrentRating := append(productIdBytes, ratingBytes...)
	// fmt.Println(" merge current rating ", mergeCurrentRating)

	// mergeCount := append(mergeCurrentRating, countBytes...)
	// fmt.Println("merge current count ", mergeCount)

	// mergeUserRating := append(mergeCount, userRatingBytes...)
	// fmt.Println("merge new user rating", mergeUserRating)

	// mergeSeller := append(sellerDIDBytes, sellerRatingBytes...)
	// fmt.Println("merge seller rating", mergeSeller)

	// mergeSellerProduct := append(mergeSeller, sellerProductCountBytes...)
	// fmt.Println("merge seller product count", mergeSellerProduct)

	// merge := append(mergeUserRating, mergeSellerProduct...)
	// fmt.Println("merge all", merge)

	// runtime := &WasmtimeRuntime{}
	// runtime.Init("rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm")
	// runtime.RunHandler(merge, int32(len(productIdBytes)), int32(len(ratingBytes)), int32(len(countBytes)), int32(len(userRatingBytes)), int32(len(sellerDIDBytes)), int32(len(sellerRatingBytes)), int32(len(sellerProductCountBytes)))
}
