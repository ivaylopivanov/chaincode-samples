package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Storage object
type Storage struct {
}

// Init will do nothing
func (s *Storage) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke Storage method
func (s *Storage) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "get":
		return get(stub, args)
	case "setPublicKey":
		return setPublicKey(stub, args)
	case "getPublicKey":
		return getPublicKey(stub, args)
	case "set":
		return set(stub, args)
	default:
		return shim.Error("Unsupported operation")
	}
}
