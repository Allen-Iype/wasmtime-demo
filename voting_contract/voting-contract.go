package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type User struct {
	DID  string
	Vote int
}

type Count struct {
	Red  int
	Blue int
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

func (r *WasmtimeRuntime) dumpOutput(pointer int32, uservote int32, red int32, blue int32, length int32) {
	fmt.Println("red :", red)
	fmt.Println("blue :", blue)
	fmt.Println("uservote :", uservote)
	r.output = make([]byte, length)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+length])

	count := Count{}
	count.Red = int(red)
	count.Blue = int(blue)

	content, err := json.Marshal(count)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("votefile.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *WasmtimeRuntime) RunHandler(data []byte, did int32, vote int32, red int32, blue int32) []byte {
	r.input = data
	r.handler.Call(r.store, did, vote, red, blue)
	fmt.Println("Result:", r.output)
	return r.output
}

func voting_contract() {

	// choices : red=1 blue =2

	randvote := rand.Intn(3-1) + 1

	newuser := User{}
	newuser.DID = "QmVkvoPGi9jvvuxsHDVJDgzPEzagBaWSZRYoRDzU244HjZ"
	newuser.Vote = randvote

	fmt.Println(" rand ", randvote)

	did := []byte(newuser.DID)
	vote := make([]byte, 4)
	binary.LittleEndian.PutUint32(vote, uint32(randvote))

	mergeuser := append(did, vote...)
	fmt.Println(" merge user ", mergeuser)

	var count Count

	byteValue, _ := ioutil.ReadFile("votefile.json")
	json.Unmarshal(byteValue, &count)

	fmt.Println("countvalue ", count)

	redvote := count.Red
	bluevote := count.Blue

	red := make([]byte, 4)
	binary.LittleEndian.PutUint32(red, uint32(redvote))

	blue := make([]byte, 4)
	binary.LittleEndian.PutUint32(blue, uint32(bluevote))

	mergevote := append(red, blue...)
	fmt.Println("mergevote ", mergevote)

	merge := append(mergeuser, mergevote...)
	fmt.Println("merge ", merge)

	runtime := &WasmtimeRuntime{}
	runtime.Init("voting_contract/target/wasm32-unknown-unknown/debug/voting_contract.wasm")
	runtime.RunHandler(merge, int32(len(did)), int32(len(vote)), int32(len(red)), int32(len(blue)))
}
