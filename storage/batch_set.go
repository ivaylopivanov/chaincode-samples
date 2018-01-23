package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type field struct {
	Key       string
	Value     string
	Signature string
}

func batchSet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	s := args[1]
	fields := []field{}

	err := json.Unmarshal([]byte(s), &fields)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	publicKey, err := keys.PublicKey(stub, alias)
	if err != nil {
		return shim.Error(err.Error())
	}

	for _, f := range fields {
		err = checkSignature(publicKey, []byte(f.Key), f.Signature)
		if err != nil {
			return shim.Error(codes.Unauthorized)
		}
		err = stub.PutState(formatNamespace(alias, f.Key), []byte(f.Value))
		if err != nil {
			return shim.Error(codes.PutState)
		}

		err = resetVerification(stub, alias, f.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

	}

	return shim.Success(nil)
}
