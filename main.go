package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"wasmtime-demo/contract"
	"wasmtime-demo/server"
)

func GenerateSmartContract() {
	did := "bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
	wasmPath := "voting_contract/target/wasm32-unknown-unknown/release/voting_contract.wasm"
	schemaPath := "store_state/vote_contract/votefile.json"
	rawCodePath := "voting_contract/src/lib.rs"
	port := "20002"
	contract.GenerateSmartContract(did, wasmPath, schemaPath, rawCodePath, port)
}

func smartContractHash() string {
	return "QmS7odhDJRG7B356PJvyXAUFiGYPPfPtzGPNHJoj6jgrgd"
}

func DeploySmartContract() {
	comment := "Deploying Test Voting Contract"
	deployerAddress := "12D3KooWCR4BW7gfPmCZhAJusqv1PoS49jgqTGvofcG4WPyg8FxV.bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
	quorumType := 2
	rbtAmount := 1
	smartContractToken := smartContractHash()
	port := "20002"
	id := contract.DeploySmartContract(comment, deployerAddress, quorumType, rbtAmount, smartContractToken, port)
	fmt.Println("Contract ID: " + id)
	contract.SignatureResponse(id, port)

}

func ExecuteSmartContractNode1() {
	comment := "Executing Test Smart Contract on Node1"
	executorAddress := "12D3KooWPfPgt1zZeQ2WcA6JGUB5u8bLNk8dPtYKiJ3MJQoeKtu7.bafybmigka7xjp73j2oy256xsmyfd66gjv6fi3ybfyjghshij4idnxui6ea"
	quorumType := 2
	smartContractData := "Red"
	smartContractToken := smartContractHash()
	port := "20009"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractNode2() {
	comment := "Executing Test Smart Contract on Node2"
	executorAddress := "12D3KooWMQaGLNGof8AfoUQh6a7aDRS2JpjkYyUrn2nzcX5bqMko.bafybmigipihqh5smgeyokgqvh7nd3yki6epaxfbefa3jxf5msw7ltj7ujm"
	quorumType := 2
	smartContractData := "Blue"
	smartContractToken := smartContractHash()
	port := "20010"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func ExecuteSmartContractNode3() {
	comment := "Executing Test Smart Contract on Node3"
	executorAddress := "12D3KooWMQuGUzoWq5EgdhBQ6YdqTQKmxpSP5s5sKyTRMsxDe1f6.bafybmidxxslkym52zhywijju54hdhvybjuf5uhj3ugcfpdr6vwko25mlma"
	quorumType := 2
	smartContractData := "Red"
	smartContractToken := smartContractHash()
	port := "20011"
	contract.ExecuteSmartContract(comment, executorAddress, quorumType, smartContractData, smartContractToken, port)
}

func SubscribeSmartContractNode1(port string) {
	contractToken := smartContractHash()
	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractNode2(port string) {
	contractToken := smartContractHash()
	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func SubscribeSmartContractNode3(port string) {
	contractToken := smartContractHash()
	contract.RegisterCallBackUrl(contractToken, "8080", "api/v1/contract-input", port)
	contract.SubscribeSmartContract(contractToken, port)
}

func main() {
	//go server.Bootup()
	fmt.Println("Server Started  2")
	go server.Bootup()

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enlighten me with the function to be executed ")
		fmt.Println(`
		1. Generate Contract 
		2. Subscribe Contract Node 1
		3. Subscribe Contract Node 2 
		4. Subscribe Contract Node 3 
		5. Deploy Contract
		6. Execute Contract Node 1 
		7. Execute Contract Node 2 
		8. Execute Contract Node 3`)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Generate Contract")
			GenerateSmartContract()
		case "2":
			fmt.Println("Subscribing Smart Contract in Node 1")
			SubscribeSmartContractNode1("20009")
		case "3":
			fmt.Println("Subscribing Smart Contract in Node 2")
			SubscribeSmartContractNode2("20010")
		case "4":
			fmt.Println("Subscribing Smart Contract in Node 3")
			SubscribeSmartContractNode3("20011")
		case "5":
			fmt.Println("Deploying Smart Contract")
			DeploySmartContract()
		case "6":
			fmt.Println("Executing Smart Contract in Node 1")
			ExecuteSmartContractNode1()
		case "7":
			fmt.Println("Executing Smart Contract in Node 2")
			ExecuteSmartContractNode2()
		case "8":
			fmt.Println("Executing Smart Contract in Node 3")
			ExecuteSmartContractNode3()
		default:
			fmt.Println("You entered an unknown number")
		}
	}

}
