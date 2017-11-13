package main

import (
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
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

const signature = "FAGwbzC323MKl8HUz+NgbOxW2NRMh7Qm+lsq5LHZt9uTYgl6sGwmO1WMD6Zh3Bu0mApwtiksncis1V5d/aJg634wZUCs4dWVj53YI1//9tRvzF/q+T3LV4DcQOFLHr9dygNUMTin/ZYqRihS6mVaF/Wg6fKNkJmom7QNA05ZxqlyjzlHNXiKQ0V6MA3XY/+PwvbGu0cveK+pBg/hHl/Rm+y3gCTT9r5VwVWoSQC88eWba2LhNPo7UR/Glr4DA/INethiS16s2KWf6gmwgvi9S/V3uH93iHQty7rVXHIfZxY6BsU+dIUcM6upGN5pQTaWbqYBKeM9WZ7DF0bBOAbVhg=="
const alias = "0ff9092fcd5908014591db39950581d8f3a28144ef7a78295e6b549e734e5ee93fb8d586c8e5969806cddbfa61ed8cf9588289dd2ef5e58cd339c0969a6d01f7"

var (
	key           = []byte("/profile/get")
	transactionID = 0
)

func TestNewChaincode(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
}

func TestSetPublicKey(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	res := setKey(stub)
	assert.Equal(t, statusOK, res.Status)
	assert.Empty(t, res.Payload)
}

func TestGetPublicKey(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	setKey(stub)
	res := getKey(stub)

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, string(res.Payload), publicKey)
}

func TestSetAndGet(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Storage))

	setKey(stub)
	value := "Some wonderful place"

	res := stub.MockInvoke(getID(), [][]byte{[]byte("set"), []byte(alias), key, []byte(value), []byte(signature)})
	assert.Equal(t, statusOK, res.Status)

	res = stub.MockInvoke(getID(), [][]byte{[]byte("get"), []byte(alias), key})

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, string(res.Payload), value)
}

func getID() string {
	transactionID = transactionID + 1
	return "TXID" + strconv.Itoa(transactionID)
}

func setKey(stub *shim.MockStub) pb.Response {
	return stub.MockInvoke(getID(), [][]byte{[]byte("setPublicKey"), []byte(alias), []byte(publicKey)})
}

func getKey(stub *shim.MockStub) pb.Response {
	return stub.MockInvoke(getID(), [][]byte{[]byte("getPublicKey"), []byte(alias)})
}
