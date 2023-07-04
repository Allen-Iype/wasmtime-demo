package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bytecodealliance/wasmtime-go"
)

type WasmtimeRuntime struct {
	store   *wasmtime.Store
	memory  *wasmtime.Memory
	handler *wasmtime.Func

	candidates []string
	votes      []int
}

func (r *WasmtimeRuntime) Init(wasmFile string, candidates []string) {
	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)
	linker.DefineWasi()
	wasiConfig := wasmtime.NewWasiConfig()
	r.store = wasmtime.NewStore(engine)
	r.store.SetWasi(wasiConfig)
	linker.FuncWrap("env", "cast_vote", r.castVote)
	linker.FuncWrap("env", "get_results", r.getResults)
	wasmBytes, _ := os.ReadFile(wasmFile)
	module, _ := wasmtime.NewModule(r.store.Engine, wasmBytes)
	instance, _ := linker.Instantiate(r.store, module)
	r.memory = instance.GetExport(r.store, "memory").Memory()
	r.handler = instance.GetFunc(r.store, "handler")
	r.candidates = candidates
	r.votes = make([]int, len(candidates))
}

func (r *WasmtimeRuntime) castVote(candidateIndex int32) {
	if candidateIndex >= 0 && int(candidateIndex) < len(r.votes) {
		r.votes[candidateIndex]++
	}
}

func (r *WasmtimeRuntime) getResults() int32 {
	maxVotes := 0
	maxIndex := -1
	for i, votes := range r.votes {
		if votes > maxVotes {
			maxVotes = votes
			maxIndex = i
		}
	}
	return int32(maxIndex)
}

func (r *WasmtimeRuntime) RunHandler() int32 {
	r.handler.Call(r.store)
	return r.getResults()
}

func main() {
	candidates := []string{"Alice", "Bob", "Charlie"}

	runtime := &WasmtimeRuntime{}
	runtime.Init("voting_contract/target/wasm32-wasi/debug/voting_contract.wasm", candidates)

	// Simulate casting votes
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			candidateIndex := int32(i % len(candidates))
			runtime.castVote(candidateIndex)
		}
	}()

	// Display results
	for {
		time.Sleep(500 * time.Millisecond)
		winnerIndex := runtime.RunHandler()
		if winnerIndex != -1 {
			fmt.Printf("Winner: %s\n", candidates[winnerIndex])
		} else {
			fmt.Println("No winner yet")
		}
	}
}
