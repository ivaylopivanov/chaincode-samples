# Wallet chaincode

## Available methods
  - get (balance for user)
  - deposit
  - withdrawal
  - history (transaction history for user)
  - getAllKeys (get all users)
  - transfer (from User A to User B)

You can have a look at the [wallet_test.go](https://github.com/ivaylopivanov/chaincode-samples/blob/master/wallet/wallet_test.go) for more information.

## Test

```sh
go test -v .
```
