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
		"503d489db69b12963afa030657c32a2fcd4547039cd6ce8d773b876b6edf1e1ea8c03ce301c60dfb822dd02a2feebbb1eae8114ffccad35b7028042d1f350879",
		`-----BEGIN PUBLIC KEY-----
		MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnll4w+0/sk3y6Zx1JTLE
		9NeioZnacvtsnwXzjN982sWiBVZBqfqHcUUc9mlR4VceGBZBgTjGstflZdK8gexD
		oLXX0Ky9HXFxSf9nAQvcZXQDS8jPiF3mQIEmwohEig4fOu0LadJ73191judS3CXJ
		WhjF+ZMVQVPRQrDEjl0vxLdp40/cz4/SkeBtBXsQcNshsSiS7HpL4XAlTUsGLfSH
		k2p7pFd3ck9bLRp2a21XczFBae1AA3gsGrKZSJ+yOkkQuJBsF5WwttE9HimPJp7u
		o2ry5mrlGaf0kEtu/erQCTTYCBmaIYBFIokHiV/UDgqQFDCLqezwXEnd974YVtwZ
		swIDAQAB
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
