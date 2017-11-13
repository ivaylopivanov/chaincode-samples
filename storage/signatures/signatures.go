package signatures

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// Verify signed message
func Verify(key, toVerify []byte, signed string) error {
	decoded, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		return err
	}

	parser, err := loadPublicKey(key)
	if err != nil {
		return err
	}

	return parser.SignVerify(toVerify, decoded)
}

//
//  All the code under this line was found on the internet
//

func loadPublicKey(key []byte) (Verifier, error) {
	return parsePublicKey(key)
}

func parsePublicKey(pemBytes []byte) (Verifier, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("Unable To Decode Public Key")
	}

	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, errors.New("Key is not supported")
	}

	return newUnsignerFromKey(rawkey)
}

// Verifier NOT A PUBLIC INTERFACE
type Verifier interface {
	SignVerify(data []byte, sig []byte) error
}

func newUnsignerFromKey(k interface{}) (Verifier, error) {
	var sshKey Verifier
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &rsaPublicKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

func (r *rsaPublicKey) SignVerify(message []byte, sig []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
}
