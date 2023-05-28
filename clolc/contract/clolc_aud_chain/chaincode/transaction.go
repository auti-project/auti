package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Transaction struct {
	ID        string `json:"id"`
	CipherRes string `json:"cipher_res"`
	CipherB   string `json:"cipher_b"`
	CipherC   string `json:"cipher_c"`
	CipherD   string `json:"cipher_d"`
}

func NewTransaction(id, cipherRes, cipherB, cipherC, cipherD string) *Transaction {
	return &Transaction{
		ID:        id,
		CipherRes: cipherRes,
		CipherB:   cipherB,
		CipherC:   cipherC,
		CipherD:   cipherD,
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
