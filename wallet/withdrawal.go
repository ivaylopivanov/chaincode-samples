package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func withdrawal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return update("withdrawal", stub, args)
}
