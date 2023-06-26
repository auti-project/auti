package organization

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

type CLOLCOrganization struct {
	ID                  TypeID
	IDHash              string
	EpochID             TypeEpochID
	epochAccumulatorMap map[[2]string]kyber.Point
	epochTXRandomness   map[[2]string]kyber.Scalar
}

func CLOLCNew(id string) *CLOLCOrganization {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	idHash := hex.EncodeToString(sha256Func.Sum(nil))
	org := &CLOLCOrganization{
		ID:                  TypeID(id),
		IDHash:              idHash,
		epochAccumulatorMap: make(map[[2]string]kyber.Point),
		epochTXRandomness:   make(map[[2]string]kyber.Scalar),
	}
	return org
}

func (c *CLOLCOrganization) SetEpochID(randID []byte) {
	c.EpochID = randID
}

func (c *CLOLCOrganization) RecordTransaction(tx *transaction.CLOLCLocalPlain) error {
	// Submit the transaction to the local chain
	sha256Func := sha256.New()
	sha256Func.Write([]byte(tx.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, randScalar, err := crypto.PedersenCommit(tx.Amount)
	if err != nil {
		return err
	}
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return err
	}
	clolcHidden := &transaction.CLOLCLocalHidden{
		CounterParty: counterPartyHash,
		Commitment:   commitmentBytes,
		Timestamp:    tx.Timestamp,
	}
	if err = c.SubmitTXLocalChain(clolcHidden); err != nil {
		return err
	}
	counterPartyHashStr := hex.EncodeToString(counterPartyHash)
	orgMapKey := IDHashKey(c.IDHash, counterPartyHashStr)
	// Accumulate the commitment to the corresponding accumulator
	if _, ok := c.epochAccumulatorMap[orgMapKey]; !ok {
		c.epochAccumulatorMap[orgMapKey] = commitment
	} else {
		c.epochAccumulatorMap[orgMapKey].Add(
			c.epochAccumulatorMap[orgMapKey],
			commitment,
		)
	}
	// Record the randomness used in the commitment
	c.epochTXRandomness[orgMapKey] = randScalar
	return nil
}

func (c *CLOLCOrganization) Accumulate(counterParty TypeID, commitment kyber.Point) {
	counterPartyHashStr := IDHashString(counterParty)
	orgMapKey := IDHashKey(c.IDHash, counterPartyHashStr)
	// Accumulate the commitment to the corresponding accumulator
	if _, ok := c.epochAccumulatorMap[orgMapKey]; !ok {
		c.epochAccumulatorMap[orgMapKey] = commitment
	} else {
		c.epochAccumulatorMap[orgMapKey].Add(
			c.epochAccumulatorMap[orgMapKey],
			commitment,
		)
	}
}

func (c *CLOLCOrganization) SubmitTXLocalChain(tx *transaction.CLOLCLocalHidden) error {
	panic("not implemented")
}

func (c *CLOLCOrganization) ComposeTXOrgChain(counterParty TypeID) (*transaction.CLOLCOrgPlain, error) {
	counterPartyHashStr := IDHashString(counterParty)
	orgMapKey := IDHashKey(c.IDHash, counterPartyHashStr)
	accumulator, ok := c.epochAccumulatorMap[orgMapKey]
	if !ok {
		return nil, fmt.Errorf("no transaction from %s to %s", c.ID, counterParty)
	}
	epochIDHashPoint := EpochIDHashPoint(c.EpochID)
	resultPoint := crypto.KyberSuite.Point().Add(accumulator, epochIDHashPoint)
	result, err := resultPoint.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return transaction.NewCLOLCOrgPlain(result), nil
}
