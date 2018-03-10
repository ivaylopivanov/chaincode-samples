package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
)

// Storage object
type Storage struct {
}

// Init will do nothing
func (s *Storage) Init(stub shim.ChaincodeStubInterface) pb.Response {
	args := []string{
		"1234567",
		`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAiRH6L4HHc8pF5ujVoNtt
fx2z9HlC+x5y3d8Tw8bkkxUZh7jn3fXL56hdgv1OAEy0vMAcZmqRs7DzSNOhY69Q
ZU0Q4gvn9Y0bYh31VaoTNMW9vvUFRkdv0bi2UW8fgY358zjd5IAwsl6tDRumSdZQ
KevKoBQJ39a4Qqf0QXcAqYzLb1iwRneP5SaNpg7ASFh0fzx6G+TR3vG9ucBnA2po
eHZu8v6C5nsx5FDsZmfS2fJW5UuAe4F+Q2+31srmBkTSVQhAdbGRxJg3zjs77R9y
aXX/GwhGUY3N7kxmez7LqBAzl3jg6mb0rFC3gWiR3UyojcYn30Hf7weVFGJuh7rm
cQIDAQAB
-----END PUBLIC KEY-----`,
		"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"}
	return create(stub, args)
}

// Invoke Storage method
func (s *Storage) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	switch fn {
	case "ping":
		return ping(stub, args)
	case "get":
		return get(stub, args)
	case "batchGet":
		return batchGet(stub, args)
	case "set":
		return set(stub, args)
	case "batchSet":
		return batchSet(stub, args)
	case "batchExchange":
		return batchExchange(stub, args)
	case "create":
		return create(stub, args)
	case "getKeys":
		return getKeys(stub, args)
	case "getPublicKey":
		return getPublicKey(stub, args)
	case "getPrivateKey":
		return getPrivateKey(stub, args)
	case "history":
		return history(stub, args)
	case "timeOfLastUpdate":
		return timeOfLastUpdate(stub, args)
	case "getDetailsForProperty":
		return getDetailsForProperty(stub, args)
	case "isVerified":
		return isVerified(stub, args)
	case "verify":
		return verify(stub, args)
	case "getVerifications":
		return getVerifications(stub, args)
	case "getVerificationFor":
		return getVerificationFor(stub, args)
	default:
		return shim.Error(codes.UnsupportedOperation)
	}
}
