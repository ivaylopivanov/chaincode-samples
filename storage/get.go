package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

func get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	key := args[1]
	ns := formatNamespace(alias, key)

	if len(args) > 2 && args[2] != "" {
		ns = formatNamespace(args[2], formatNamespace(alias, key))
	}

	value, err := stub.GetState(ns)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	if len(value) == 0 {
		return shim.Error(codes.NotFound)
	}

	return shim.Success(value)
}
