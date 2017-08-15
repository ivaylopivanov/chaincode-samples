package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error("Not enough arguments")
	}
	from := args[0]
	to := args[1]
	value := args[2]

	currentFrom, err := stub.GetState(from)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance for 'from' error: %s", err))
	}

	newBalanceFrom, err := calculateNewBalance("withdrawal", value, currentFrom)
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance 'from' error: %s", err))
	}

	currentTo, err := stub.GetState(to)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance for 'to' error: %s", err))
	}

	newBalanceTo, err := calculateNewBalance("deposit", value, currentTo)
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance 'to' error: %s", err))
	}

	if err := stub.PutState(from, newBalanceFrom); err != nil {
		return shim.Error(fmt.Sprintf("Save balance 'from' error: %s", err))
	}

	if err := stub.PutState(to, newBalanceTo); err != nil {
		return shim.Error(fmt.Sprintf("Save balance 'to' error: %s", err))
	}

	return shim.Success(newBalanceFrom)
}
