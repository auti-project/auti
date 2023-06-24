package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type CLOLCOrgPlain struct {
	Accumulator []byte
}

func NewCLOLCOrgPlain(accumulator []byte) *CLOLCOrgPlain {
	return &CLOLCOrgPlain{
		Accumulator: accumulator,
	}
}

func (c *CLOLCOrgPlain) ToOnChain() *CLOLCOrgOnChain {
	accumulatorString := hex.EncodeToString(c.Accumulator)
	return NewCLOLCOrgOnChain(accumulatorString)
}

type CLOLCOrgOnChain struct {
	Accumulator string `json:"accumulator"`
}

func NewCLOLCOrgOnChain(accumulator string) *CLOLCOrgOnChain {
	return &CLOLCOrgOnChain{
		Accumulator: accumulator,
	}
}

func (c *CLOLCOrgOnChain) ToPlain() (*CLOLCOrgPlain, error) {
	accumulatorBytes, err := hex.DecodeString(c.Accumulator)
	if err != nil {
		return nil, err
	}
	return NewCLOLCOrgPlain(accumulatorBytes), nil
}

func (c *CLOLCOrgOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(c)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
