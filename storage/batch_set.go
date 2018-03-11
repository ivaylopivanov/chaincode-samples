package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type field struct {
	Property  string
	Hash      string
	Signature string
}

func batchSet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]
	s := args[1]
	fields := []field{}

	err := json.Unmarshal([]byte(s), &fields)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	publicKey, err := keys.PublicKey(stub, id)
	if err != nil {
		return shim.Error(err.Error())
	}

	for _, f := range fields {
		err = checkSignature(publicKey, []byte(f.Property), f.Signature)
		if err != nil {
			return shim.Error(codes.Unauthorized)
		}
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
