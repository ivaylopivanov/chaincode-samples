package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 4 {
		return shim.Error("Not enough arguments")
	}
	from := args[0]
	to := args[1]
	currency := args[2]
	value := args[3]

	currentFrom, err := stub.GetState(from)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance for 'from' error: %s", err))
	}

	stateFrom, err := newState(currentFrom)
	if err != nil {
		return shim.Error(fmt.Sprintf("Parse current balance for 'from' error: %s", err))
	}

	newBalanceFrom, err := calculateNewBalance("withdrawal", value, stateFrom.get(currency))
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance for 'from' error: %s", err))
	}

	stateFrom.set(currency, newBalanceFrom)

	currentTo, err := stub.GetState(to)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance for 'to' error: %s", err))
	}

	stateTo, err := newState(currentTo)
	if err != nil {
		return shim.Error(fmt.Sprintf("Parse current balance for 'to' error: %s", err))
	}

	newBalanceTo, err := calculateNewBalance("deposit", value, stateTo.get(currency))
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance 'to' error: %s", err))
	}

	stateTo.set(currency, newBalanceTo)

	b, err := stateFrom.toJSON()
	if err != nil {
		return shim.Error(fmt.Sprintf("toJSON for 'from' error: %s", err))
	}

	if err := stub.PutState(from, b); err != nil {
		return shim.Error(fmt.Sprintf("Save balance 'from' error: %s", err))
	}

	b, err = stateTo.toJSON()
	if err != nil {
		return shim.Error(fmt.Sprintf("toJSON for 'to' error: %s", err))
	}

	if err := stub.PutState(to, b); err != nil {
		return shim.Error(fmt.Sprintf("Save balance 'to' error: %s", err))
	}

	return shim.Success(nil)
}
