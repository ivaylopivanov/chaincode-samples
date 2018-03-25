package main

import (
	"encoding/json"

	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (s Storage) batchExchange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	from := args[0]
	to := args[1]
	f := args[2]
	fields := []field{}

	err := json.Unmarshal([]byte(f), &fields)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	publicKey, err := keys.PublicKey(stub, from)
	if err != nil {
		return shim.Error(err.Error())
	}

	for _, f := range fields {
		err = checkSignature(publicKey, []byte(formatNamespace(to, f.Property)), f.Signature)
		err = resetVerificationFor(stub, from, to, f.Property)
		if err != nil {
			return shim.Error(err.Error())
		}
		if err != nil {
			return shim.Error(codes.Unauthorized)
		}
		err = stub.PutState(formatNamespace(from, formatNamespace(to, f.Property)), []byte(f.Hash))
		if err != nil {
			return shim.Error(codes.PutState)
		}

	}

	return shim.Success(nil)
}
