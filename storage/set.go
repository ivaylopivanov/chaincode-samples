package main

import (
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const statusOK = int32(200)

func set(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 4 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	key := args[1]
	value := args[2]
	signature := args[3]

	err := checkIdentity(stub, alias, key, signature)
	if err != nil {
		return shim.Error(codes.Unauthorized)
	}

	err = stub.PutState(formatNamespace(alias, key), []byte(value))
	if err != nil {
		return shim.Error(codes.PutState)
	}

	err = resetVerification(stub, alias, key)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
