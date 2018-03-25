package signatures

import (
	"crypto"
	"crypto/rand"
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

// Sign message
func Sign(key, data []byte) ([]byte, error) {
	parser, err := loadPrivateKey(key)
	if err != nil {
		return nil, err
	}

	s, err := parser.Sign(data)
	if err != nil {
		return nil, err
	}

	return []byte(base64.StdEncoding.EncodeToString(s)), nil
}

//
//  All the code under this line was found on the internet
//

func loadPublicKey(key []byte) (Verifier, error) {
	return parsePublicKey(key)
}

func parsePublicKey(b []byte) (Verifier, error) {
	block, _ := pem.Decode(b)
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

// loadPrivateKey parses a PEM encoded private key.
func loadPrivateKey(b []byte) (Signer, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}
