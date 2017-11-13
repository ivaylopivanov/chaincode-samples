package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ivaylopivanov/chaincode-samples/multi-currency-wallet/numbers"
	"github.com/ivaylopivanov/chaincode-samples/multi-currency-wallet/rsa"
	"github.com/ivaylopivanov/chaincode-samples/multi-currency-wallet/transactions"
)

func transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return shim.Error("Not enough arguments")
	}

	if !isAuthorized(stub, args) {
		return shim.Error("Not authorized")
	}

	from := args[0]
	to := args[1]
	currency := args[2]
	value := args[3]

	currentFrom, err := stub.GetState(from)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance for 'from' error: %s", err))
	}

	stateFrom, err := newState(currentFrom)
	if err != nil {
		return shim.Error(fmt.Sprintf("Parse current balance for 'from' error: %s", err))
	}

	newBalanceFrom, err := calculateNewBalance("withdrawal", value, stateFrom.get(currency))
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance for 'from' error: %s", err))
	}

	stateFrom.set(currency, newBalanceFrom)

	currentTo, err := stub.GetState(to)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get current balance for 'to' error: %s", err))
	}

	stateTo, err := newState(currentTo)
	if err != nil {
		return shim.Error(fmt.Sprintf("Parse current balance for 'to' error: %s", err))
	}

	newBalanceTo, err := calculateNewBalance("deposit", value, stateTo.get(currency))
	if err != nil {
		return shim.Error(fmt.Sprintf("Update balance 'to' error: %s", err))
	}

	stateTo.set(currency, newBalanceTo)

	b, err := stateFrom.toJSON()
	if err != nil {
		return shim.Error(fmt.Sprintf("toJSON for 'from' error: %s", err))
	}

	if err := stub.PutState(from, b); err != nil {
		return shim.Error(fmt.Sprintf("Save balance 'from' error: %s", err))
	}

	b, err = stateTo.toJSON()
	if err != nil {
		return shim.Error(fmt.Sprintf("toJSON for 'to' error: %s", err))
	}

	if err := stub.PutState(to, b); err != nil {
		return shim.Error(fmt.Sprintf("Save balance 'to' error: %s", err))
	}

	t, err := transactions.Save(stub, from, to, currency, value)

	if err != nil {
		return shim.Error(fmt.Sprintf("Save transaction error: %s", err))
	}

	return shim.Success(t)
}

func isAuthorized(stub shim.ChaincodeStubInterface, args []string) bool {
	from := args[0]
	to := args[1]
	currency := args[2]
	value := args[3]
	signature := args[4]

	attempt, err := incrementAttempts(stub, from)
	if err != nil {
		return false
	}

	label := from + to + currency + value + attempt

	key, err := getKeyFor(stub, from)
	if err != nil {
		return false
	}

	err = rsa.VerifySignature(key, label, signature)
	if err != nil {
		return false
	}

	return true
}

func getTransferAttempts(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Not enough arguments")
	}
	key := "attempts-" + args[0]
	current, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get attempts error: %s", err))
	}
	return shim.Success(current)
}

func incrementAttempts(stub shim.ChaincodeStubInterface, user string) (string, error) {
	key := "attempts-" + user
	current, err := stub.GetState(key)
	if err != nil {
		return "", err
	}
	attempts := int64(0)
	if len(current) > 0 {
		attempts = numbers.ToInt64(current)
	}
	newAttempt := numbers.Int64ToString(attempts + 1)
	err = stub.PutState(key, []byte(newAttempt))
	return newAttempt, nil
}
