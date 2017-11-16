package main

import (
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	key := args[0]
	value := args[1]

	if key == publicKeyNamespace {
		return shim.Error(codes.BadRequest)
	}

	v, err := stub.GetState(key)
	if err != nil || len(v) > 0 {
		return shim.Error(codes.AlreadyExists)
	}

	err = stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}
