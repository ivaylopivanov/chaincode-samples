package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]

	v, err := stub.GetState(alias)
	if err != nil || len(v) > 0 {
		return shim.Error(codes.AlreadyExists)
	}

	k := keys{
		Public:  args[1],
		Private: args[2],
	}

	b, err := json.Marshal(k)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	err = stub.PutState(alias, b)
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}
