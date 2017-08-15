package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getAllKeys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	iter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("GetAllKeys state error: %s", err))
	}
	defer iter.Close()

	var keys []string

	for iter.HasNext() {

		res, err := iter.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("GetAllKeys iteration error: %s", err))
		}

		keys = append(keys, res.Key)
	}

	result, err := json.Marshal(keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetAllKeys marshal error: %s", err))
	}

	return shim.Success(result)
}
