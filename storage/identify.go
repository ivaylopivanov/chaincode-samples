package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
)

var identifycationKey = "idkey"

func (s Storage) identify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error(codes.WrongAmountOfArguments)
	}

	id := args[0]
	signature := args[1]
	hash := args[2]

	err := identify(stub, id, signature, hash)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func identify(stub shim.ChaincodeStubInterface, id, signature, value string) error {
	publicKey, err := keys.PublicKey(stub, id)
	if err != nil {
		return errors.New(codes.GetState)
	}

	s := strings.Split(signature, "-tp-")
	if len(s) != 2 {
		return errors.New(codes.Unknown)
	}

	sig := s[0]
	timestamp := s[1]
	key := sig + "-" + timestamp
	if cacheCheck(id, key) {
		return errors.New(codes.KeyViolation)
	}

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if len(timestamp) == 13 {
		i = i / 1000
	}
	if err != nil {
		return errors.New(codes.BadRequest)
	}
	then := time.Unix(i, 0)
	if time.Since(then).Hours() > 5 {
		return errors.New(codes.TimeViolation)
	}

	sha := sha256.New()
	sha.Write([]byte(value))
	hash := hex.EncodeToString(sha.Sum(nil))
	signed := hash + "-" + timestamp

	err = checkSignature(publicKey, []byte(signed), sig)
	if err != nil {
		return errors.New(signed)
	}

	return nil
}
