package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var publicKeyNamespace = "key-"

func setPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Not enough arguments")
	}
	key := publicKeyNamespace + args[0]
	publicKey := publicKeyNamespace + args[1]
	err := stub.PutState(key, []byte(publicKey))
	if err != nil {
		return shim.Error(fmt.Sprintf("Set public key error: %s", err))
	}
	return shim.Success(nil)
}

func getPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Not enough arguments")
	}
	key := publicKeyNamespace + args[0]
	current, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get public key error: %s", err))
	}
	return shim.Success(current)
}

func getKeyFor(stub shim.ChaincodeStubInterface, user string) ([]byte, error) {
	key := publicKeyNamespace + user
	publicKey, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}
