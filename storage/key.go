package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/signatures"
)

var publicKeyNamespace = "key"

func setPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	res := getPublicKey(stub, args)
	if len(res.Payload) > 0 {
		return shim.Error(codes.AlreadyExists)
	}

	alias := args[0]
	publicKey := args[1]

	err := stub.PutState(formatNamespace(alias, publicKeyNamespace), []byte(publicKey))
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}

func getPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.NotEnoughArguments)
	}

	current, err := getKeyFor(stub, args[0])

	if err != nil {
		return shim.Error(codes.GetState)
	}

	return shim.Success(current)
}

func verify(stub shim.ChaincodeStubInterface, alias, key, signature string) error {
	publicKey, err := getKeyFor(stub, alias)
	if err != nil {
		return errors.New(codes.NotFound)
	}

	return signatures.Verify(publicKey, []byte(key), signature)
}

func getKeyFor(stub shim.ChaincodeStubInterface, alias string) ([]byte, error) {
	publicKey, err := stub.GetState(formatNamespace(alias, publicKeyNamespace))

	if err != nil {
		return nil, err
	}

	if len(publicKey) == 0 {
		return nil, errors.New(codes.NotFound)
	}

	return publicKey, nil
}
