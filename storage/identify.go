package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/multi-currency-wallet/numbers"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
)

var identifycationKey = "idkey"

type identifyRes struct {
	Success bool
	Current int64
}

func (s Storage) identify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]
	signature := args[1]

	publicKey, err := keys.PublicKey(stub, id)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	key := formatNamespace(id, identifycationKey)

	b, err := stub.GetState(key)

	if err != nil {
		return shim.Error(codes.GetState)
	}

	n := int64(0)

	if b != nil {
		n = numbers.ByteToInt64(b)
	}
	n = n + 1

	res := identifyRes{}
	res.Current = n

	err = checkSignature(publicKey, []byte(numbers.Int64ToString(n)), signature)
	if err == nil {
		res.Success = true
		err = stub.PutState(key, numbers.Int64ToByte(n))
		if err != nil {
			return shim.Error(codes.PutState)
		}
	}

	b, err = json.Marshal(&res)
	if err != nil {
		return shim.Error(codes.Unknown)
	}

	return shim.Success(b)
}
