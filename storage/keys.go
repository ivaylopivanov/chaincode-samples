package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
	"github.com/ivaylopivanov/chaincode-samples/storage/signatures"
)

func (s Storage) getKeys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]

	b, err := keys.Get(stub, id)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func (s Storage) getPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]

	b, err := keys.PublicKey(stub, id)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func (s Storage) getPrivateKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]

	b, err := keys.PrivateKey(stub, id)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(b)
}

func checkIdentity(stub shim.ChaincodeStubInterface, id, signed, signature string) error {
	publicKey, err := keys.PublicKey(stub, id)
	if err != nil {
		return err
	}
	return checkSignature(publicKey, []byte(signed), signature)
}

func checkSignature(publicKey, signed []byte, signature string) error {
	return signatures.Verify(publicKey, signed, signature)
}
