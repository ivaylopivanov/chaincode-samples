package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func update(fn string, stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Not enough arguments")
	}

	key := args[0]
	value := args[1]

	current, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance error: %s", err))
	}

	newBalance, err := calculateNewBalance(fn, value, current)
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance error: %s", err))
	}

	if err := stub.PutState(key, newBalance); err != nil {
		return shim.Error(fmt.Sprintf("Save balance error: %s", err))
	}

	return shim.Success(newBalance)
}

func calculateNewBalance(fn string, change string, current []byte) ([]byte, error) {
	value, err := stringToFloat64(change)
	if err != nil {
		return nil, err
	}

	if value < 0 {
		// we cannot have negative values which may produce an undesirable result
		// eg: 20 - -10 or 20 + - 10
		return nil, errors.New("Passing a negative value")
	}

	if current == nil {
		// the user has no balance yet - eg new user
		return float64ToByte(value), nil
	}

	currentBalance := byteToFloat64(current)
	if fn == "deposit" {
		return float64ToByte(currentBalance + value), nil
	}
	// otherwhise it's a withdrawal
	newBalance := currentBalance - value
	if newBalance < 0 {
		return nil, errors.New("insufficient funds")
	}

	return float64ToByte(newBalance), nil
}
