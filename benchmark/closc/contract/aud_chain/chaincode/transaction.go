package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	Commitment string `json:"commitment"`
	BatchProof string `json:"batch_proof"`
}

func NewTransaction(commitment, batchProof string) *Transaction {
	return &Transaction{
		Commitment: commitment,
		BatchProof: batchProof,
	}
}

func (t *Transaction) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(t)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
