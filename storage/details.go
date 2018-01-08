package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

type details struct {
	Time          int64
	Value         string
	Verifications []verification
}

func getDetailsForKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	from := args[0]
	to := args[1]
	key := args[2]

	value, err := stub.GetState(formatNamespace(from, formatNamespace(to, key)))
	if err != nil {
		return shim.Error(codes.GetState)
	}

	if len(value) == 0 {
		return shim.Error(codes.NotFound)
	}

	d := details{}
	d.Value = string(value)

	iter, err := stub.GetHistoryForKey(formatNamespace(from, key))
	if err != nil {
		return shim.Error(codes.GetState)
	}
	defer iter.Close()

	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return shim.Error(codes.GetHistory)
		}

		d.Time = res.Timestamp.Seconds
	}

	v, err := fetchVerifications(stub, from, key)
	if err != nil {
		return shim.Error(err.Error())
	}

	d.Verifications = v

	b, err := json.Marshal(d)
	if err != nil {
		return shim.Error(codes.Unknown)
	}

	return shim.Success(b)
}
