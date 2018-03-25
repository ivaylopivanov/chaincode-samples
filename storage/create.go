package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (s Storage) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]

	v, err := stub.GetState(id)
	if err != nil || len(v) > 0 {
		return shim.Error(codes.AlreadyExists)
	}

	k := keys.Keys{
		Public:  args[1],
		Private: args[2],
	}

	b, err := json.Marshal(k)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	err = stub.PutState(id, b)
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}
