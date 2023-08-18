package main

import (
	"fmt"
	"wasmtime-demo/contract"
	"wasmtime-demo/server"
)

/*From the discussions what I understood is since the Rust part or the Smart contract is written by the DAPP, they would know the structure
of the data which is being passed or the data which is being used, so this structure which the smart contract deployer or the dapp has
created must be known to the Go-package. So this structure which is being written in the Rust code must be provided to the GO-code
as the Schema  */

func main() {
	//go server.Bootup()
	fmt.Println("Server Started  2")

	comment := "Wasm_test"
	//deployerAddress := "12D3KooWCR4BW7gfPmCZhAJusqv1PoS49jgqTGvofcG4WPyg8FxV.bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"
	smartContractToken := "QmPa3rqRjUThHtzH57RwBTXvU6K5EmqRTWmKbdnFoWgE1w"
	executorAddress := "bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa"

	//contract.GenerateSmartContract("bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa", "rating_contract/target/wasm32-unknown-unknown/release/rating_contract.wasm", "store_state/rating_contract/rating.json", "rating_contract/src/lib.rs", "20002")

	//QmPa3rqRjUThHtzH57RwBTXvU6K5EmqRTWmKbdnFoWgE1w
	//QmPa3rqRjUThHtzH57RwBTXvU6K5EmqRTWmKbdnFoWgE1w

	// id := contract.DeploySmartContract(comment, deployerAddress, 2, 1, smartContractToken, "20002")
	// fmt.Println("Smart Contract Deployed with ID: ", id)
	// contract.SignatureResponse(id, "20002")
	// contractToken := "QmPa3rqRjUThHtzH57RwBTXvU6K5EmqRTWmKbdnFoWgE1w"
	port := "20002"
	// contract.SubscribeSmartContract(contractToken, port)
	contract.ExecuteSmartContract(comment, executorAddress, 2, "5", smartContractToken, port)

	//contract.DeploySmartContract("Test","bafybmifb4rbwykckpbcnekcha23nckrldhkcqyrhegl7oz44njgci5vhqa",)
	//deploy
	//execute
	//subscribe should be called by others
	server.Bootup()

}
