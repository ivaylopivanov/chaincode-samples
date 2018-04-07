package main

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

func (s Storage) batchGet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]
	properties := strings.Split(args[1], ",")

	var wg sync.WaitGroup

	results := map[string]string{}

	for _, pr := range properties {
		wg.Add(1)

		ns := formatNamespace(id, pr)

		if len(args) > 2 && args[2] != "" {
			ns = formatNamespace(args[2], formatNamespace(id, pr))
		}

		value, err := stub.GetState(ns)
		if err == nil && len(value) > 0 {
			results[pr] = string(value)
		}

		wg.Done()
	}

	wg.Wait()

	b, err := json.Marshal(results)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	return shim.Success(b)
}
