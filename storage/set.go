package main

import (
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const statusOK = int32(200)

func (s Storage) set(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 4 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]
	signature := args[1]
	property := args[2]
	hash := args[3]
	err := identify(stub, id, signature, hash)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = resetVerification(stub, id, property)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(formatNamespace(id, property), []byte(hash))
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}
