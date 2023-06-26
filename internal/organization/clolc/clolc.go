package clolc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"go.dedis.ch/kyber/v3"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/organization"
	"github.com/auti-project/auti/internal/transaction/clolc"
)

type Organization struct {
	ID                  organization.TypeID
	IDHash              string
	EpochID             organization.TypeEpochID
	epochAccumulatorMap map[[2]string]kyber.Point
	epochTXRandomness   map[[2]string]kyber.Scalar
}

func New(id string) *Organization {
	defer organization.Sha256Func.Reset()
	organization.Sha256Func.Write([]byte(id))
	idHash := hex.EncodeToString(organization.Sha256Func.Sum(nil))
	org := &Organization{
		ID:                  organization.TypeID(id),
		IDHash:              idHash,
		epochAccumulatorMap: make(map[[2]string]kyber.Point),
		epochTXRandomness:   make(map[[2]string]kyber.Scalar),
	}
	return org
}

func (c *Organization) SetEpochID(randID []byte) {
	c.EpochID = randID
}

func (c *Organization) RecordTransaction(tx *clolc.LocalPlain) error {
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
	clolcHidden := &clolc.LocalHidden{
		CounterParty: counterPartyHash,
		Commitment:   commitmentBytes,
		Timestamp:    tx.Timestamp,
	}
	if err = c.SubmitTXLocalChain(clolcHidden); err != nil {
		return err
	}
	counterPartyHashStr := hex.EncodeToString(counterPartyHash)
	orgMapKey := organization.IDHashKey(c.IDHash, counterPartyHashStr)
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

func (c *Organization) Accumulate(counterParty organization.TypeID, commitment kyber.Point) {
	counterPartyHashStr := organization.IDHashString(counterParty)
	orgMapKey := organization.IDHashKey(c.IDHash, counterPartyHashStr)
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

func (c *Organization) SubmitTXLocalChain(tx *clolc.LocalHidden) error {
	panic("not implemented")
}

func (c *Organization) ComposeTXOrgChain(counterParty organization.TypeID) (*clolc.OrgPlain, error) {
	counterPartyHashStr := organization.IDHashString(counterParty)
	orgMapKey := organization.IDHashKey(c.IDHash, counterPartyHashStr)
	accumulator, ok := c.epochAccumulatorMap[orgMapKey]
	if !ok {
		return nil, fmt.Errorf("no transaction from %s to %s", c.ID, counterParty)
	}
	epochIDHashPoint := organization.EpochIDHashPoint(c.EpochID)
	resultPoint := crypto.KyberSuite.Point().Add(accumulator, epochIDHashPoint)
	result, err := resultPoint.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return clolc.NewOrgPlain(result), nil
}
