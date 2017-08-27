package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Transaction for each update on the balance
type Transaction struct {
	Timestamp int64
	Value     []byte
	ID        string
}

func history(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Not enough arguments")
	}

	user := args[0]
	iter, err := stub.GetHistoryForKey(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("History state error: %s", err))
	}
	defer iter.Close()

	var keys []Transaction

	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("History iteration error: %s", err))
		}

		keys = append(keys, Transaction{
			ID:        res.TxId,
			Timestamp: res.Timestamp.Seconds,
			Value:     res.Value,
		})
	}

	result, err := json.Marshal(keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("History marshal error: %s", err))
	}

	return shim.Success(result)
}
