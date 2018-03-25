package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

// Transaction for each update on the balance
type Transaction struct {
	Timestamp int64
	Value     []byte
	ID        string
}

func (s Storage) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]
	key := args[1]

	iter, err := stub.GetHistoryForKey(formatNamespace(id, key))
	if err != nil {
		return shim.Error(codes.GetState)
	}
	defer iter.Close()

	var keys []Transaction

	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return shim.Error(codes.GetHistory)
		}

		keys = append(keys, Transaction{
			ID:        res.TxId,
			Timestamp: res.Timestamp.Seconds,
			Value:     res.Value,
		})
	}

	result, err := json.Marshal(keys)
	if err != nil {
		return shim.Error(codes.Unknown)
	}

	return shim.Success(result)
}
