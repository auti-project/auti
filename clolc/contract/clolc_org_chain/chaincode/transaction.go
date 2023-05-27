package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	Accumulator string `json:"accumulator"`
}

func NewTransaction(accumulator string) *Transaction {
	return &Transaction{
		Accumulator: accumulator,
	}
}

func (c *Transaction) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(c)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
