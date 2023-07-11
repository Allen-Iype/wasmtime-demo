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

	input  []byte
	output []byte
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

func (r *WasmtimeRuntime) dumpOutput(pointer int32, userRating int32, ratingCount int32, length int32) {
	r.output = make([]byte, length)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+length])

	review := ProductReview{}
	review.DID = string(r.output)
	review.Rating = int(userRating)
	review.RatingCount = int(ratingCount)

	content, err := json.Marshal(review)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("store_state/rating_contract/rating.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *WasmtimeRuntime) RunHandler(data []byte, did int32, rating int32, count int32) []byte {
	r.input = data
	r.handler.Call(r.store, did, rating, count)
	fmt.Println("Result:", r.output)
	return r.output
}

func main() {

	var review ProductReview
	stateUpdate, _ := os.ReadFile("store_state/rating_contract/rating.json")
	json.Unmarshal(stateUpdate, &review)
	randomRating := rand.Intn(5) + 1
	ratingCount := review.RatingCount

	fmt.Println("DID : ", review.DID)
	fmt.Println("Rating : ", randomRating)
	fmt.Println("RatingCount : ", ratingCount)

	did := []byte(review.DID)
	rating := make([]byte, 4)
	binary.LittleEndian.PutUint32(rating, uint32(randomRating))
	count := make([]byte, 4)
	binary.LittleEndian.PutUint32(count, uint32(ratingCount))
	mergeuser := append(did, rating...)
	fmt.Println(" merge user ", mergeuser)
	merge := append(mergeuser, count...)
	fmt.Println("merge ", merge)

	runtime := &WasmtimeRuntime{}
	runtime.Init("voting_contract/target/wasm32-unknown-unknown/debug/voting_contract.wasm")
	runtime.RunHandler(merge, int32(len(did)), int32(len(rating)), int32(len(count)))
}
