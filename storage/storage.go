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
		"934da5b412e08c45f026d540f00c0514ae0d1bcaa05b6da3c601e56bcf9174394912e53fc6926a920c3f49c537aa7cd7392632f35c85082919d2cbba310f3763",
		`-----BEGIN PUBLIC KEY-----
		MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1wHRIOtBsCyv5UQZuigr
		vmE+kg5ai1fNC1wvC9OhwiZJdARyrX3VgJE3f17z12tklTnYYKYmFF15qHhjoRl2
		aLNCWOd9aKw7dAsZRZcnWK/dJRo8ocNST9xssBAUy94M+1k0H+BJBRIrMJxshiZX
		qPIUIZtg0Q5tCnn4nAeENSkdRPrvZfKImbY97N5pbpzF2QphiAsFlvjAOTC9cEtG
		YH/qUyZXfQ/NMAfwaqs+fEbWnUGGUSYTeNioj9xm5eNrXzCnqFP+6o2Z6xgwmxCG
		9Tm84JOR1V4vY0bX7NqYr8fvBK7YboUGi7AGRW397IG1sTUWEmii7p/EkdXA8A9n
		twIDAQAB
		-----END PUBLIC KEY-----`,
		"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",
	}
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
	case "create":
		return create(stub, args)
	case "getKeys":
		return getKeys(stub, args)
	case "getPublicKey":
		return getPublicKey(stub, args)
	case "getPrivateKey":
		return getPrivateKey(stub, args)
	default:
		return shim.Error(codes.UnsupportedOperation)
	}
}
