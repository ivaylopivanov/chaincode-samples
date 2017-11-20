package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/signatures"
)

type keys struct {
	Public  string
	Private string
}

func getKeys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]

	b, err := stub.GetState(alias)

	if err != nil {
		return shim.Error(codes.GetState)
	}

	return shim.Success(b)
}

func verify(stub shim.ChaincodeStubInterface, alias, key, signature string) error {
	b, err := stub.GetState(alias)
	if err != nil {
		return errors.New(codes.GetState)
	}

	k := &keys{}
	json.Unmarshal(b, k)

	if err != nil {
		return errors.New(codes.NotFound)
	}

	return signatures.Verify([]byte(k.Public), []byte(key), signature)
}
