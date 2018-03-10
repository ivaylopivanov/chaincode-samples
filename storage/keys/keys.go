package keys

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

// Keys obj
type Keys struct {
	Public  string
	Private string
}

// Get keys
func Get(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	b, err := stub.GetState(id)

	if err != nil {
		return nil, errors.New(codes.GetState)
	}

	return b, nil
}

// PublicKey for id
func PublicKey(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	b, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New(codes.GetState)
	}

	k := &Keys{}
	json.Unmarshal(b, k)

	if err != nil {
		return nil, errors.New(codes.NotFound)
	}

	return []byte(k.Public), nil
}

// PrivateKey for id
func PrivateKey(stub shim.ChaincodeStubInterface, id string) ([]byte, error) {
	b, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New(codes.GetState)
	}

	k := &Keys{}
	json.Unmarshal(b, k)

	if err != nil {
		return nil, errors.New(codes.NotFound)
	}

	return []byte(k.Private), nil
}
