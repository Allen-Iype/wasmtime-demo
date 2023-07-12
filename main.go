package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/bytecodealliance/wasmtime-go"
)

type WasmtimeRuntime struct {
	store   *wasmtime.Store
	memory  *wasmtime.Memory
	handler *wasmtime.Func

	input     []byte
	output    []byte
	outputPtr []byte
}

type ProductReview struct {
	DID         string
	Rating      int
	RatingCount int
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

func (r *WasmtimeRuntime) dumpOutput(pointer int32, userRating int32, ratingCount int32, length int32, outputPtr int32, outputLen int32) {
	r.output = make([]byte, length)
	r.outputPtr = make([]byte, outputLen)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+length])
	copy(r.outputPtr, r.memory.UnsafeData(r.store)[pointer:outputPtr+outputLen])
	review := ProductReview{}
	review.DID = string(r.output)
	review.Rating = int(userRating)
	review.RatingCount = int(ratingCount)
	fmt.Println("OutputPtr :", string(r.outputPtr))
	content, err := json.Marshal(review)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("store_state/rating_contract/rating.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *WasmtimeRuntime) RunHandler(data []byte, didLength int32, ratingLength int32, countLength int32, userRatingLength int32) []byte {
	r.input = data
	r.handler.Call(r.store, didLength, ratingLength, countLength, userRatingLength)
	fmt.Println("Result:", r.output)
	return r.output
}

func main() {

	var review ProductReview
	stateUpdate, _ := os.ReadFile("store_state/rating_contract/rating.json")
	json.Unmarshal(stateUpdate, &review)
	currentRating := review.Rating
	ratingCount := review.RatingCount

	fmt.Println("DID : ", review.DID)
	fmt.Println("Current Rating : ", currentRating)
	fmt.Println("Current Rating Count : ", ratingCount)

	randomRating := rand.Intn(5) + 1 //A random rating from 1-5 given for testing[Here it is considered as the rating a user gave]

	did := []byte(review.DID)
	//the current average rating of the product
	rating := make([]byte, 4)
	binary.LittleEndian.PutUint32(rating, uint32(currentRating))
	//the total count of the ratings received
	count := make([]byte, 4)
	binary.LittleEndian.PutUint32(count, uint32(ratingCount))
	//the new rating given by the user
	userRating := make([]byte, 4)
	binary.LittleEndian.PutUint32(userRating, uint32(randomRating))

	mergeCurrentRating := append(did, rating...)
	fmt.Println(" merge current rating ", mergeCurrentRating)

	mergeCount := append(mergeCurrentRating, count...)
	fmt.Println("merge current count ", mergeCount)

	merge := append(mergeCount, userRating...)
	fmt.Println("merge new user rating", merge)

	runtime := &WasmtimeRuntime{}
	runtime.Init("rating_contract/target/wasm32-unknown-unknown/debug/rating_contract.wasm")
	runtime.RunHandler(merge, int32(len(did)), int32(len(rating)), int32(len(count)), int32(len(userRating)))
}
