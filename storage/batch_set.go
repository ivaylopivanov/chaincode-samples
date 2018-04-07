package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type field struct {
	Property string `json:"property"`
	Hash     string `json:"hash"`
}

func (s Storage) batchSet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]
	signature := args[1]
	f := args[2]
	fields := []field{}

	err := identify(stub, id, signature, f)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal([]byte(f), &fields)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	for _, f := range fields {
		err = stub.PutState(formatNamespace(id, f.Property), []byte(f.Hash))
		if err != nil {
			return shim.Error(codes.PutState)
		}

		err = resetVerification(stub, id, f.Property)
		if err != nil {
			return shim.Error(err.Error())
		}

	}

	return shim.Success(nil)
}
