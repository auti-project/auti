package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"go.dedis.ch/kyber/v3"
)

type AudPlain struct {
	Commitment []byte
}

func NewAudPlain(commitment []byte) *AudPlain {
	return &AudPlain{
		Commitment: commitment,
	}
}

func NewAudPlainFromPoint(commitment kyber.Point) (*AudPlain, error) {
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return NewAudPlain(commitmentBytes), nil
}

func (a *AudPlain) ToOnChain() *AudOnChain {
	return NewAudOnChain(hex.EncodeToString(a.Commitment))
}

type AudOnChain struct {
	Commitment string `json:"commitment"`
}

func NewAudOnChain(commitment string) *AudOnChain {
	return &AudOnChain{
		Commitment: commitment,
	}
}

func (a *AudOnChain) ToPlain() (*AudPlain, error) {
	commitment, err := hex.DecodeString(a.Commitment)
	if err != nil {
		return nil, err
	}
	return NewAudPlain(commitment), nil
}

func (a *AudOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(a)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
