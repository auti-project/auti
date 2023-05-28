package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type CLOLCAudPlain struct {
	ID        []byte
	CipherRes []byte
	CipherB   []byte
	CipherC   []byte
	CipherD   []byte
}

func NewCLOLCAudPlain(id, cipherRes, cipherB, cipherC, cipherD []byte) *CLOLCAudPlain {
	return &CLOLCAudPlain{
		ID:        id,
		CipherRes: cipherRes,
		CipherB:   cipherB,
		CipherC:   cipherC,
		CipherD:   cipherD,
	}
}

func (c *CLOLCAudPlain) ToOnChain() *CLOLCAudOnChain {
	return &CLOLCAudOnChain{
		ID:        hex.EncodeToString(c.ID),
		CipherRes: hex.EncodeToString(c.CipherRes),
		CipherB:   hex.EncodeToString(c.CipherB),
		CipherC:   hex.EncodeToString(c.CipherC),
		CipherD:   hex.EncodeToString(c.CipherD),
	}
}

type CLOLCAudOnChain struct {
	ID        string `json:"id"`
	CipherRes string `json:"cipher_res"`
	CipherB   string `json:"cipher_b"`
	CipherC   string `json:"cipher_c"`
	CipherD   string `json:"cipher_d"`
}

func (c *CLOLCAudOnChain) ToPlain() (*CLOLCAudPlain, error) {
	id, err := hex.DecodeString(c.ID)
	if err != nil {
		return nil, err
	}
	cipherRes, err := hex.DecodeString(c.CipherRes)
	if err != nil {
		return nil, err
	}
	cipherB, err := hex.DecodeString(c.CipherB)
	if err != nil {
		return nil, err
	}
	cipherC, err := hex.DecodeString(c.CipherC)
	if err != nil {
		return nil, err
	}
	cipherD, err := hex.DecodeString(c.CipherD)
	if err != nil {
		return nil, err
	}
	return NewCLOLCAudPlain(id, cipherRes, cipherB, cipherC, cipherD), nil
}

func (c *CLOLCAudOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(c)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
