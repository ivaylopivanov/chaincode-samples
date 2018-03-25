package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
)

type verification struct {
	UserID    int64
	Signature string
	Status    string
	Timestamp string
}

var (
	// StatusRejected for verification
	StatusRejected = "rejected"
)

func (s Storage) getVerifications(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]
	property := args[1]

	ns := formatVerificationNamespace(id, property)

	b, err := stub.GetState(ns)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	return shim.Success(b)
}

func (s Storage) getVerificationFor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]
	property := args[1]
	toCheckFor := args[2]

	ns := formatVerificationNamespace(id, property)

	b, err := stub.GetState(ns)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	verifications := []verification{}
	err = json.Unmarshal([]byte(b), &verifications)
	if err != nil {
		return shim.Error(codes.BadRequest)
	}

	idToCheckFor, _ := stringToInt64(toCheckFor)

	for _, v := range verifications {
		if v.UserID == idToCheckFor {
			b, err = json.Marshal(&v)
			if err != nil {
				return shim.Error(codes.Unknown)
			}
			return shim.Success(b)
		}
	}

	return shim.Error(codes.NotFound)
}

func (s Storage) isVerified(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(codes.NotEnoughArguments)
	}

	id := args[0]
	property := args[1]

	ns := formatVerificationNamespace(id, property)

	b, err := stub.GetState(ns)
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
		return shim.Error(codes.NotVerified)
	}

	for _, v := range verifications {
		if v.Status == StatusRejected {
			return shim.Error(codes.NotVerified)
		}
	}

	return shim.Success(nil)
}

func (s Storage) verify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 6 {
		return shim.Error(codes.NotEnoughArguments)
	}

	from := args[0]
	to := args[1]
	property := args[2]
	signature := args[3]
	status := args[4]
	timestamp := args[5]

	publicKey, err := keys.PublicKey(stub, from)
	if err != nil {
		return shim.Error(codes.GetState)
	}

	ns := formatVerificationNamespace(to, property)

	err = checkSignature(publicKey, []byte(ns), signature)
	if err != nil {
		return shim.Error(codes.Unauthorized)
	}

	b, err := stub.GetState(ns)
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

	id, _ := stringToInt64(from)

	for k, v := range verifications {
		if v.UserID == id {
			verifications = append(verifications[:k], verifications[k+1:]...)
			break
		}
	}

	verifications = append(verifications, verification{
		UserID:    id,
		Signature: signature,
		Status:    status,
		Timestamp: timestamp,
	})

	b, err = json.Marshal(&verifications)
	if err != nil {
		return shim.Error(codes.Unknown)
	}

	err = stub.PutState(ns, b)
	if err != nil {
		return shim.Error(codes.PutState)
	}

	return shim.Success(nil)
}

func resetVerification(stub shim.ChaincodeStubInterface, id, property string) error {
	err := stub.PutState(formatVerificationNamespace(id, property), []byte(nil))
	if err != nil {
		return errors.New(codes.PutState)
	}
	return nil
}

func resetVerificationFor(stub shim.ChaincodeStubInterface, from, to, property string) error {
	ns := formatVerificationNamespace(from, property)

	b, err := stub.GetState(ns)
	if err != nil {
		return errors.New(codes.GetState)
	}

	vers := []verification{}
	if len(b) > 0 {
		err = json.Unmarshal(b, &vers)
	}

	id, _ := stringToInt64(to)

	for k, v := range vers {
		if v.UserID == id {
			vers = append(vers[:k], vers[k+1:]...)
		}
	}

	b, err = json.Marshal(&vers)
	if err != nil {
		return errors.New(codes.Unknown)
	}

	err = stub.PutState(ns, b)
	if err != nil {
		return errors.New(codes.PutState)
	}

	return err
}

func fetchVerifications(stub shim.ChaincodeStubInterface, id, property string) ([]verification, error) {
	ns := formatVerificationNamespace(id, property)

	b, err := stub.GetState(ns)
	if err != nil {
		return nil, errors.New(codes.GetState)
	}

	v := []verification{}
	if len(b) > 0 {
		err = json.Unmarshal(b, &v)
	}

	return v, err
}

func formatVerificationNamespace(id, property string) string {
	return id + "-verified" + property
}

func stringToInt64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
