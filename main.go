package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	UserId     string
	BillAmount int
	Balance    int
	UserStatus string
}

//	type PaymentStatus struct {
//		paymentStatus string
//	}
func (r *WasmtimeRuntime) loadInput(pointer int32) {
	copy(r.memory.UnsafeData(r.store)[pointer:pointer+int32(len(r.input))], r.input)
}

func (r *WasmtimeRuntime) dumpOutput(pointer int32, billPaid int32, length int32) {
	fmt.Println("billPaid :", billPaid)
	r.output = make([]byte, length)
	copy(r.output, r.memory.UnsafeData(r.store)[pointer:pointer+length])
	user := User{}
	if billPaid == 0 {
		fmt.Println("bill not paid, user has insufficient balance deactivating account")
		user.UserId = string(r.output)
		user.UserStatus = "Inactive"
		user.Balance = 20
		user.BillAmount = 30

	} else {
		fmt.Println("bill paid, user has sufficient balance activating account")
		user.UserId = string(r.output)
		user.UserStatus = "Active"
		user.Balance = 10
		user.BillAmount = 0
	}
	fmt.Println(user)
	content, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("userbill.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *WasmtimeRuntime) RunHandler(data []byte, userid int32, billAmount int32, balance int32, userStatus int32) []byte {
	r.input = data
	r.handler.Call(r.store, userid, billAmount, balance, userStatus)
	fmt.Println("Result:", r.output)
	return r.output
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

func main() {

	// choices : red=1 blue =2

	//randvote := rand.Intn(3-1) + 1

	newuser := User{}
	newuser.UserId = "User1"
	newuser.BillAmount = 30
	newuser.Balance = 20
	newuser.UserStatus = "Active"

	userid := []byte(newuser.UserId)
	fmt.Println("UserId input :", userid)
	billAmount := make([]byte, 4)
	binary.LittleEndian.PutUint32(billAmount, uint32(newuser.BillAmount))
	balance := make([]byte, 4)
	binary.LittleEndian.PutUint32(balance, uint32(newuser.Balance))
	userStatus := []byte(newuser.UserStatus)
	fmt.Println("UserStatus input :", userStatus)

	//since wasm has a linear memory we are appending all the bytes linearly
	//mergeuser := append(userid, billAmount, balance, userStatus)
	mergeuser := append(userid, billAmount...)
	mergeuser = append(mergeuser, balance...)
	mergeuser = append(mergeuser, userStatus...)
	fmt.Println(" merge user ", mergeuser)

	runtime := &WasmtimeRuntime{}
	runtime.Init("bill_contract/target/wasm32-unknown-unknown/debug/bill_contract.wasm")
	runtime.RunHandler(mergeuser, int32(len(userid)), int32(len(billAmount)), int32(len(balance)), int32(len(userStatus)))
}
