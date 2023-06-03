package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/auti-project/auti/internal/crypto"
	"go.dedis.ch/kyber/v3"
)

const (
	amountAmplifier = 100
)

type CLOLCLocalPlain struct {
	CounterParty string
	Amount       int64
	Timestamp    int64
}

func NewCLOLCLocalPlain(counterParty string, amount float64, timestamp int64) *CLOLCLocalPlain {
	amountInt := int64(amount * amountAmplifier)
	return &CLOLCLocalPlain{
		CounterParty: counterParty,
		Amount:       amountInt,
		Timestamp:    timestamp,
	}
}

func NewPairCLOLCLocalPlain(
	fromID, toID string,
	amount float64,
	timestamp int64,
) (*CLOLCLocalPlain, *CLOLCLocalPlain) {
	amountInt := int64(amount * amountAmplifier)
	return &CLOLCLocalPlain{
			CounterParty: toID,
			Amount:       amountInt,
			Timestamp:    timestamp,
		}, &CLOLCLocalPlain{
			CounterParty: fromID,
			Amount:       -amountInt,
			Timestamp:    timestamp,
		}
}

func (c *CLOLCLocalPlain) Hide() (*CLOLCLocalHidden, kyber.Point, kyber.Scalar, error) {
	sha256Func := sha256.New()
	sha256Func.Write([]byte(c.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, randScalar, err := crypto.PedersenCommit(c.Amount)
	if err != nil {
		return nil, nil, nil, err
	}
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return nil, nil, nil, err
	}
	return NewCLOLCLocalHidden(
		counterPartyHash,
		commitmentBytes,
		c.Timestamp,
	), commitment, randScalar, nil
}

type CLOLCLocalHidden struct {
	CounterParty []byte
	Commitment   []byte
	Timestamp    int64
}

func NewCLOLCLocalHidden(counterParty, commitment []byte, timestamp int64) *CLOLCLocalHidden {
	return &CLOLCLocalHidden{
		CounterParty: counterParty,
		Commitment:   commitment,
		Timestamp:    timestamp,
	}
}

func (c *CLOLCLocalHidden) ToOnChain() *CLOLCLocalOnChain {
	timestampStr := strconv.FormatInt(c.Timestamp, 10)
	return NewCLOLCLocalOnChain(
		hex.EncodeToString(c.CounterParty),
		hex.EncodeToString(c.Commitment),
		timestampStr,
	)
}

type CLOLCLocalOnChain struct {
	CounterParty string `json:"counter_party"`
	Commitment   string `json:"commitment"`
	Timestamp    string `json:"timestamp"`
}

func NewCLOLCLocalOnChain(counterParty, commitment, timestamp string) *CLOLCLocalOnChain {
	return &CLOLCLocalOnChain{
		CounterParty: counterParty,
		Commitment:   commitment,
		Timestamp:    timestamp,
	}
}

func (c *CLOLCLocalOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(c)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}

func (c *CLOLCLocalOnChain) ToHidden() (*CLOLCLocalHidden, error) {
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
	return NewCLOLCLocalHidden(counterParty, commitment, timestamp), nil
}
