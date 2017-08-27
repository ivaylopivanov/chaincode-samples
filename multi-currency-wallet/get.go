package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Not enough arguments")
	}

	user := args[0]

	current, err := stub.GetState(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get balance error: %s", err))
	}

	state, err := newState(current)
	if err != nil {
		return shim.Error(fmt.Sprintf("Parse current balance error: %s", err))
	}

	b, err := state.toJSON()
	if err != nil {
		return shim.Error(fmt.Sprintf("toJSON error: %s", err))
	}

	return shim.Success(b)
}
