package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/multi-currency-wallet/transactions"
)

func history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Not enough arguments")
	}

	user := args[0]
	b, err := transactions.Get(stub, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("History state error: %s", err))
	}

	return shim.Success(b)
}
