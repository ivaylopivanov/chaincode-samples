package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func batchExchange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	from := args[0]
	to := args[1]
	s := args[2]
	fields := []field{}

	err := json.Unmarshal([]byte(s), &fields)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	publicKey, err := keys.PublicKey(stub, from)
	if err != nil {
		return shim.Error(err.Error())
	}

	for _, f := range fields {
		err = verifySignature(publicKey, []byte(formatNamespace(to, f.Key)), f.Signature)
		if err != nil {
			return shim.Error(codes.Unauthorized)
		}
		err = stub.PutState(formatNamespace(from, formatNamespace(to, f.Key)), []byte(f.Value))
		if err != nil {
			return shim.Error(codes.PutState)
		}

	}

	return shim.Success(nil)
}
