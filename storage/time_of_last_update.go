package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

func (s Storage) timeOfLastUpdate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]
	property := args[1]

	iter, err := stub.GetHistoryForKey(formatNamespace(id, property))
	if err != nil {
		return shim.Error(codes.GetState)
	}
	defer iter.Close()

	var time int64

	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return shim.Error(codes.GetHistory)
		}

		time = res.Timestamp.Seconds
	}

	return shim.Success([]byte(strconv.FormatInt(time, 10)))
}
