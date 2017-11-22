package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
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
	case "ping":
		return ping(stub, args)
	case "get":
		return get(stub, args)
	case "multipleGet":
		return multipleGet(stub, args)
	case "set":
		return set(stub, args)
	case "create":
		return create(stub, args)
	case "getKeys":
		return getKeys(stub, args)
	default:
		return shim.Error(codes.UnsupportedOperation)
	}
}
