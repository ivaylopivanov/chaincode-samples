package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type wallet struct {
}

// Init will do nothing
func (w *wallet) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke chaincode method
func (w *wallet) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "deposit":
		return deposit(stub, args)
	case "withdrawal":
		return withdrawal(stub, args)
	case "transfer":
		return transfer(stub, args)
	case "get":
		return get(stub, args)
	case "getAllKeys":
		return getAllKeys(stub, args)
	case "history":
		return history(stub, args)
	default:
		return shim.Error("Unsupported operation")
	}
}
