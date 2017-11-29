package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
	"github.com/ivaylopivanov/chaincode-samples/storage/signatures"
)

// type keys struct {
// 	Public  string
// 	Private string
// }

func getKeys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]

	b, err := keys.Get(stub, alias)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func getPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]

	b, err := keys.PublicKey(stub, alias)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func getPrivateKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]

	b, err := keys.PrivateKey(stub, alias)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func verify(stub shim.ChaincodeStubInterface, alias, key, signature string) error {
	publicKey, err := keys.PublicKey(stub, alias)
	if err != nil {
		return err
	}
	return verifySignature(publicKey, []byte(key), signature)
}

func verifySignature(publicKey, key []byte, signature string) error {
	return signatures.Verify(publicKey, key, signature)
}
