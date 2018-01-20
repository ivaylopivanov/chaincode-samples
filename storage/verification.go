package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
)

type verification struct {
	Alias     string
	Signature string
	Status    string
	Timestamp string
}

func getVerifications(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	key := args[1]

	verificationKey := formatVerificationNamespace(alias, key)

	b, err := stub.GetState(verificationKey)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	return shim.Success(b)
}

func getVerificationFor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	key := args[1]
	aliasToCheckFor := args[2]

	verificationKey := formatVerificationNamespace(alias, key)

	b, err := stub.GetState(verificationKey)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	verifications := []verification{}
	err = json.Unmarshal([]byte(b), &verifications)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	for _, v := range verifications {
		if v.Alias == aliasToCheckFor {
			b, err = json.Marshal(&v)
			if err != nil {
				return shim.Error(codes.Unknown)
			}
			return shim.Success(b)
		}
	}

	return shim.Error(codes.NotFound)
}

func isVerified(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	alias := args[0]
	key := args[1]

	verificationKey := formatVerificationNamespace(alias, key)

	b, err := stub.GetState(verificationKey)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	if b == nil {
		return shim.Error(codes.NotFound)
	}

	verifications := []verification{}
	err = json.Unmarshal(b, &verifications)
	if err != nil {
		return shim.Error(codes.Unknown)
	}

	if len(verifications) == 0 {
		return shim.Error(codes.NotFound)
	}

	return shim.Success(nil)
}

func verify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 6 {
		return shim.Error(codes.NotEnoughArguments)
	}

	from := args[0]
	to := args[1]
	key := args[2]
	signature := args[3]
	status := args[4]
	timestamp := args[5]

	publicKey, err := keys.PublicKey(stub, from)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	verificationKey := formatVerificationNamespace(to, key)

	err = checkSignature(publicKey, []byte(verificationKey), signature)
	if err != nil {
		return shim.Error(codes.Unauthorized)
	}

	b, err := stub.GetState(verificationKey)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	verifications := []verification{}

	if len(b) > 0 {
		err = json.Unmarshal([]byte(b), &verifications)
		if err != nil {
			return shim.Error(codes.Unknown)
		}
	}

	for k, v := range verifications {
		if v.Alias == from {
			verifications = append(verifications[:k], verifications[k+1:]...)
			break
		}
	}

	verifications = append(verifications, verification{
		Alias:     from,
		Signature: signature,
		Status:    status,
		Timestamp: timestamp,
	})

	b, err = json.Marshal(&verifications)
	if err != nil {
		return shim.Error(codes.Unknown)
	}

	err = stub.PutState(verificationKey, b)
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}

func resetVerification(stub shim.ChaincodeStubInterface, alias, key string) error {
	err := stub.PutState(formatVerificationNamespace(alias, key), []byte(nil))
	if err != nil {
		return errors.New(codes.PutState)
	}
	return nil
}

func fetchVerifications(stub shim.ChaincodeStubInterface, alias, key string) ([]verification, error) {
	verificationKey := formatVerificationNamespace(alias, key)

	b, err := stub.GetState(verificationKey)
	if err != nil {
		return nil, errors.New(codes.GetState)
	}

	v := []verification{}
	if len(b) > 0 {
		err = json.Unmarshal(b, &v)
	}

	return v, err
}

func formatVerificationNamespace(alias, key string) string {
	return alias + "-verified" + key
}
