package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(new(Wallet))
	if err != nil {
		fmt.Printf("Starting chaincode error: %s", err)
	}
}
