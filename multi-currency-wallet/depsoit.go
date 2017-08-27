package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func deposit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return update("deposit", stub, args)
}
