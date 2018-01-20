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
		"e2128074995f7d2cdc84a1961f6f8be2d5ee182e043af7b96b6cc78bcfcafc47520d496f0c74f3475c43b1289507c03e38261f42806df98035f6c7d59aa7f8d9",
		`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlPQhXavQr2EL4+Pao0iV
DNYJPTWeZ1swBVX7T5pcZHU1OVeWrdDWGK8iyk7HHMpOAZmDgi1gRcuVjakNXdGr
oYwTs4btLcS/ZoLnphheTwvr8DPgtcfl+iuSbEgDsOoLDoYGpOTGMiPbmRJrR5ja
BCSdQts0dWhQk4hPZapv6AA8RUy1ZREAR2U/K0H31Gkvx9v61kJduFzz4DkBs2nW
l4zsqr2E0mv9UzQDgY8rpl0nL8Jo2Ut25tkJ3BOnwfY7v4DXNlawxd3UJJ2VeGVx
kvmEyG7tv46dLFV7rp1kyWENRBvo0sQZ0rD07RoNRicx0ARABFZac9Gwa1hB+119
XQIDAQAB
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
	case "getDetailsForKey":
		return getDetailsForKey(stub, args)
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
