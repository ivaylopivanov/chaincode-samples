package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/storage/codes"
	"github.com/ivaylopivanov/chaincode-samples/storage/keys"
	"github.com/ivaylopivanov/chaincode-samples/storage/signatures"
	"github.com/stretchr/testify/assert"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgoYU8KfRqknnscTuWRcu
7Y5fRD8287L46HuYdUwfVgECgIYIzdki05ODbHK2dXNsPrM1FyW8kay2zD/VXJeL
969vycO0LAE/GjKTKYWatqKCiesJ53xNoNnb/W3/Zoo0D/fzkkw5HR1a3mG3DlSJ
egoDud+i/Wwfvg9T6YXWwbz0LD7SZd9qw80nMlHu/PL4KA2DQ/AHS0W0Y1mEOim9
8I4qVNwTMSHpmjV3oxl4tZITLd50ch9OZmcUuKHl1DAva4rQ4m27PXZqFaHMNoBm
T0HNw9+hA7tQrqoirZQ0RciOHk42/Sn3ymh1YaRBjtbVosp7cnLVeUjEfJiUHs7U
VQIDAQAB
-----END PUBLIC KEY-----`

// just for reference - not used at the moment
const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAgoYU8KfRqknnscTuWRcu7Y5fRD8287L46HuYdUwfVgECgIYI
zdki05ODbHK2dXNsPrM1FyW8kay2zD/VXJeL969vycO0LAE/GjKTKYWatqKCiesJ
53xNoNnb/W3/Zoo0D/fzkkw5HR1a3mG3DlSJegoDud+i/Wwfvg9T6YXWwbz0LD7S
Zd9qw80nMlHu/PL4KA2DQ/AHS0W0Y1mEOim98I4qVNwTMSHpmjV3oxl4tZITLd50
ch9OZmcUuKHl1DAva4rQ4m27PXZqFaHMNoBmT0HNw9+hA7tQrqoirZQ0RciOHk42
/Sn3ymh1YaRBjtbVosp7cnLVeUjEfJiUHs7UVQIDAQABAoIBAHJbn8UqhBy7G/E5
JcuQ8GDauMVGzZK/YC3w/CbpRxtHTzXkOZqBgG33dNJzv0Ewm8pjoURin9DSjmZu
FzZE4TFl2H/io91aSjtdzGo40NDrmYvVDpxu4GTp/EETOw5QUEUdbZ7kgbXsnkzx
OD2p+7mdRJ56PofjT/xp2Y9k4EBlT1x1870El0WupSr9MMI7BfPEUd8qWVy8UFIW
KIbyA1V41JEE2hkiWxZxSps1ZbaFOsWd3PCkLPcnXcLhDag2zJA38P8w6hqXqseI
PkX6tFB5NgoV9Pr9nhsaou7Tmy5xs2rv8VcfUDkCfaP5QfjmB1rsLMHzvoAw05Jk
cxgjbAECgYEA6ZSlWyn3NwZIJsk5rVGVhE+NFbaK40jN9hMxlnsvvGDauCOv878Q
XFWbMT3PcrSbHEhE6YXxC2do+l0LzGz3AuRkVDW/lDWYCkvLNcR8tRjgcqp5Io2G
AfxR3UCZORCgWlpjbNhssvtVL0i8x4YzyaPe4JaYrCKksnOLNt9vkUECgYEAjw00
olwqmE+oFdI/mS5V1jiB1FjSJ5AFzOh1/fzGtO7cHX4U+pe7HKDDtyMJo/ALoLHh
ONLjSFA9lg8M2aSM/Oq2anParU0paVU6smqOP5S7B4RGztnZrIiKLZf3NZ7+tpI8
YG4/Aq5JAx8Y65sRXSJN2eQJGOVv0b/BovTDahUCgYAtwYaa0x+wUbS0lFqODxtA
7exvQnD2kP53o87k8YCqYDa4N5VyJA1qaQKbpMYMbECuS6HkNO7BEyLHWI7FHttM
X70fmd/Lgqwj6DEIeVMMjrD5BVfxYtPLc8f9lXfua8ldqbMsUUEJ1p4bQx5n32wp
pcY1LIr/vVGR+3xb5W0PgQKBgET23dgkTNivFl6mxMhpgTJMfbLMu0wdb95wd2ni
Mj3KJc4GGcER40AS9SfWOXCSalinSOgnKzGSlY5BZTWL0figgx7hCZyg5YUFdM3M
9xUJ8/zUtXpFpl46WTtP1vs/0pZb+8WVgMUfFs3tcQss7/sRbazM9eHNwtHA+24/
R/JdAoGAE6NXAFMfpWqBFztQ8hEhv67T2KiJF9zwYC7dbyUwEYuKgexTKi7l72/q
nXbrJYAG5mOjW//wwfRz3hMlnzKDIvZ946dHIPsguF2qRYjmNndQmAoShohnWlrT
371aa5Wj7/Gv+DamJItScMs9hnZhdk4SffnIvrsAcZu08VNLXLc=
-----END RSA PRIVATE KEY-----`

const id = "1234567"

var (
	property      = []byte("/profile/get")
	transactionID = 0
)

func TestNewChaincode(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
}

func TestPing(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	res := stub.MockInvoke(getID(), [][]byte{[]byte("ping")})
	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, "pong", string(res.Payload))
}

func TestCreate(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	res := mockCreate(stub)
	assert.Equal(t, statusOK, res.Status)

	res = mockCreate(stub)
	assert.Equal(t, codes.AlreadyExists, res.Message)
}

func TestGetKeys(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	mockCreate(stub)
	res := getMockKeys(stub)

	k := &keys.Keys{}
	json.Unmarshal(res.Payload, k)

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, publicKey, k.Public)
	assert.Equal(t, privateKey, k.Private)
}

func TestSetAndGet(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	mockCreate(stub)
	value := "Some wonderful place 1"

	res := mockSet(stub, property, []byte(value))
	assert.Equal(t, statusOK, res.Status)

	res = stub.MockInvoke(getID(), [][]byte{[]byte("get"), []byte(id), property})

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, string(res.Payload), value)
}

func TestBatchGetWithSingleKey(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	mockCreate(stub)
	value := "Some wonderful place 2"

	expected := `{"/profile/get":"Some wonderful place 2"}`

	res := mockSet(stub, property, []byte(value))
	assert.Equal(t, statusOK, res.Status)

	res = stub.MockInvoke(getID(), [][]byte{[]byte("batchGet"), []byte(id), property})

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, expected, string(res.Payload))
}

func TestBatchSet(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	mockCreate(stub)
	value := "Some wonderful place"
	h := sha256.New()
	h.Write([]byte(value))
	hash := hex.EncodeToString(h.Sum(nil))

	expected := `{"/profile/get":"` + hash + `"}`

	f := []field{
		field{
			Property: string(property),
			Hash:     hash,
		},
		field{
			Property: string(property),
			Hash:     hash,
		},
	}

	b, _ := json.Marshal(f)

	s := sign(b)

	res := stub.MockInvoke(getID(), [][]byte{[]byte("batchSet"), []byte(id), s, b})
	assert.Equal(t, statusOK, res.Status)

	res = stub.MockInvoke(getID(), [][]byte{[]byte("batchGet"), []byte(id), property})

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, expected, string(res.Payload))
	stub.MockInvoke(getID(), [][]byte{[]byte("history"), []byte(id), property})
}

func TestIdentify(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	mockCreate(stub)

	sha := sha256.New()
	sha.Write([]byte("some random hash"))
	hash := hex.EncodeToString(sha.Sum(nil))

	s := sign([]byte(hash))

	res := stub.MockInvoke(getID(), [][]byte{[]byte("identify"), []byte(id), s, []byte(hash)})
	assert.Equal(t, statusOK, res.Status)
}

func getID() string {
	transactionID = transactionID + 1
	return "TXID" + strconv.Itoa(transactionID)
}

func mockCreate(stub *shim.MockStub) pb.Response {
	return stub.MockInvoke(getID(), [][]byte{[]byte("create"), []byte(id), []byte(publicKey), []byte(privateKey)})
}

func mockSet(stub *shim.MockStub, k, v []byte) pb.Response {
	s := sign(v)
	return stub.MockInvoke(getID(), [][]byte{[]byte("set"), []byte(id), s, k, v})
}

func getMockKeys(stub *shim.MockStub) pb.Response {
	return stub.MockInvoke(getID(), [][]byte{[]byte("getKeys"), []byte(id)})
}

func sign(v []byte) []byte {
	t := time.Now().Format("20060102150405")
	sha := sha256.New()
	sha.Write(v)
	value := string(hex.EncodeToString(sha.Sum(nil))) + "-" + t
	signed, _ := signatures.Sign([]byte(privateKey), []byte(value))
	s := string(signed) + "-tp-" + t
	return []byte(s)
}
