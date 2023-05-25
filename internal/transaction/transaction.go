package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/auti-project/auti/internal/crypto"
)

const (
	amountAmplifier = 100
)

type CLOLCPlain struct {
	CounterParty string
	Amount       int64
	Timestamp    int64
}

func NewCLOLCPlain(counterParty string, amount float64, timestamp int64) *CLOLCPlain {
	amountInt := int64(amount * amountAmplifier)
	return &CLOLCPlain{
		CounterParty: counterParty,
		Amount:       amountInt,
		Timestamp:    timestamp,
	}
}

func (c *CLOLCPlain) Hide() (*CLOLCHidden, error) {
	sha256Func := sha256.New()
	sha256Func.Write([]byte(c.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, err := crypto.PedersenCommit(c.Amount)
	if err != nil {
		return nil, err
	}
	return NewCLOLCHidden(
		counterPartyHash,
		commitment,
		c.Timestamp,
	), nil
}

type CLOLCHidden struct {
	CounterParty []byte
	Commitment   []byte
	Timestamp    int64
}

func NewCLOLCHidden(counterParty, commitment []byte, timestamp int64) *CLOLCHidden {
	return &CLOLCHidden{
		CounterParty: counterParty,
		Commitment:   commitment,
	}
}

func (c *CLOLCHidden) ToOnChain() *CLOLCOnChain {
	timestampStr := strconv.FormatInt(c.Timestamp, 10)
	return NewCLOLCOnChain(
		hex.EncodeToString(c.CounterParty),
		hex.EncodeToString(c.Commitment),
		timestampStr,
	)
}

type CLOLCOnChain struct {
	CounterParty string `json:"counter_party"`
	Commitment   string `json:"commitment"`
	Timestamp    string `json:"timestamp"`
}

func NewCLOLCOnChain(counterParty, commitment, timestamp string) *CLOLCOnChain {
	return &CLOLCOnChain{
		CounterParty: counterParty,
		Commitment:   commitment,
		Timestamp:    timestamp,
	}
}

func (c *CLOLCOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(c)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}

func (c *CLOLCOnChain) ToHidden() (*CLOLCHidden, error) {
	counterParty, err := hex.DecodeString(c.CounterParty)
	if err != nil {
		return nil, err
	}
	commitment, err := hex.DecodeString(c.Commitment)
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(c.Timestamp, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewCLOLCHidden(counterParty, commitment, timestamp), nil
}
