package main

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
)

var (
	userA             = []byte("A")
	depositAmount     = "100"
	bDepositAmount    = []byte(depositAmount)
	withdrawalAmount  = "11.115111"
	bWithdrawalAmount = []byte(withdrawalAmount)
	transactionID     = 0
	statusOK          = int32(200)
	statusError       = int32(500)
)

func TestNewChaincode(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
}

func TestGet(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))

	res := stub.MockInvoke(getID(), [][]byte{[]byte("get"), userA})

	assert.Equal(t, statusOK, res.Status)
	assert.Empty(t, res.Payload)
}

func TestDeposit(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))

	res := makeDeposit(stub)

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, toFloat64(depositAmount), byteToFloat64(res.Payload))
}

func TestWithdrawal(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))

	makeDeposit(stub)

	res := makeWithdrawal(stub)

	newBalance := toFloat64(depositAmount) - toFloat64(withdrawalAmount)
	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, newBalance, byteToFloat64(res.Payload))

	// just to make sure
	res = makeWithdrawal(stub)
	newBalance = newBalance - toFloat64(withdrawalAmount)

	assert.Equal(t, statusOK, res.Status)
	assert.Equal(t, newBalance, byteToFloat64(res.Payload))
}

func TestHistory(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))

	res := stub.MockInvoke(getID(), [][]byte{[]byte("history"), userA})

	// The GetHistoryForKey is not implemented by the mockstub, yet
	assert.Equal(t, statusError, res.Status)
	assert.Equal(t, "History state error: Not Implemented", res.Message)
}

func TestGetAllKeys(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))

	stub.MockInvoke(getID(), [][]byte{[]byte("deposit"), []byte("B"), bDepositAmount})
	stub.MockInvoke(getID(), [][]byte{[]byte("deposit"), []byte("C"), bDepositAmount})
	stub.MockInvoke(getID(), [][]byte{[]byte("deposit"), []byte("C"), bDepositAmount})
	stub.MockInvoke(getID(), [][]byte{[]byte("deposit"), []byte("A"), bDepositAmount})

	res := stub.MockInvoke(getID(), [][]byte{[]byte("getAllKeys")})

	result := []string{}
	err := json.Unmarshal(res.Payload, &result)

	expected := []string{"A", "B", "C"}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, statusOK, res.Status)
}

func TestTransfer(t *testing.T) {
	stub := shim.NewMockStub("mockStub", new(Wallet))

	makeDeposit(stub)

	userB := []byte("B")

	res := stub.MockInvoke(getID(), [][]byte{[]byte("transfer"), userA, userB, bWithdrawalAmount})
	assert.Equal(t, statusOK, res.Status)

	res = stub.MockInvoke(getID(), [][]byte{[]byte("get"), userA})

	expectedBalanceForA := toFloat64(depositAmount) - toFloat64(withdrawalAmount)
	assert.Equal(t, expectedBalanceForA, byteToFloat64(res.Payload))

	res = stub.MockInvoke(getID(), [][]byte{[]byte("get"), userB})

	expectedBalanceForB := toFloat64(withdrawalAmount)
	assert.Equal(t, expectedBalanceForB, byteToFloat64(res.Payload))
}

func makeDeposit(stub *shim.MockStub) pb.Response {
	return stub.MockInvoke(getID(), [][]byte{[]byte("deposit"), userA, bDepositAmount})
}

func makeWithdrawal(stub *shim.MockStub) pb.Response {
	return stub.MockInvoke(getID(), [][]byte{[]byte("withdrawal"), userA, bWithdrawalAmount})
}

func getID() string {
	transactionID = transactionID + 1
	return "TXID" + strconv.Itoa(transactionID)
}

func toFloat64(number string) float64 {
	n, _ := stringToFloat64(number)
	return n
}
