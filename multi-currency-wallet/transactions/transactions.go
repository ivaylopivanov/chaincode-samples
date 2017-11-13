package transactions

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Transaction object
type Transaction struct {
	ID       string `json:"id" msgpack:"id"`
	From     string `json:"from" msgpack:"from"`
	To       string `json:"to" msgpack:"to"`
	Currency string `json:"currency" msgpack:"currency"`
	Value    string `json:"value" msgpack:"value"`
	Date     int64  `json:"date" msgpack:"date"`
}

var namespace = "transactions-"

// Save transaction
func Save(stub shim.ChaincodeStubInterface, from, to, currency, value string) ([]byte, error) {
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, err
	}

	t := &Transaction{
		ID:       stub.GetTxID(),
		From:     from,
		To:       to,
		Currency: currency,
		Value:    value,
		Date:     timestamp.GetSeconds(),
	}
	return t.save(stub)
}

// Get transactions for given user
func Get(stub shim.ChaincodeStubInterface, user string) ([]byte, error) {
	return stub.GetState(namespace + user)
}

func (t *Transaction) save(stub shim.ChaincodeStubInterface) ([]byte, error) {

	t.Date = t.Date * 1000

	stateFrom, err := Get(stub, t.From)
	if err != nil {
		return nil, err
	}

	stateTo, err := Get(stub, t.To)
	if err != nil {
		return nil, err
	}

	previousTransactionsFrom := []Transaction{}
	if stateFrom != nil {
		err = json.Unmarshal(stateFrom, &previousTransactionsFrom)
		if err != nil {
			return nil, err
		}
	}

	previousTransactionsTo := []Transaction{}
	if stateTo != nil {
		err = json.Unmarshal(stateTo, &previousTransactionsTo)
		if err != nil {
			return nil, err
		}
	}

	previousTransactionsFrom = append(previousTransactionsFrom, *t)
	previousTransactionsTo = append(previousTransactionsTo, *t)

	byteFrom, err := json.Marshal(previousTransactionsFrom)
	if err != nil {
		return nil, err
	}

	byteTo, err := json.Marshal(previousTransactionsTo)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(namespace+t.From, byteFrom)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(namespace+t.To, byteTo)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return b, nil
}
