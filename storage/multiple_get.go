package main

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

func multipleGet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	keys := strings.Split(args[1], ",")

	var wg sync.WaitGroup

	results := map[string]string{}

	for _, k := range keys {
		wg.Add(1)

		value, err := stub.GetState(formatNamespace(alias, k))
		if err == nil && len(value) > 0 {
			results[k] = string(value)
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
