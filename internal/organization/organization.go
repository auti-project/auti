package organization

import (
	"crypto/sha256"
	"math/big"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

type TypeID string

type Organization struct {
	ID                  TypeID
	epochRandID         *big.Int
	epochAccumulatorMap map[[2]TypeID]kyber.Point
}

func New(id string) *Organization {
	org := &Organization{
		ID:                  TypeID(id),
		epochAccumulatorMap: make(map[[2]TypeID]kyber.Point),
	}
	return org
}

func (o *Organization) SetEpochID(randID *big.Int) {
	o.epochRandID = randID
}

func (o *Organization) RecordTransaction(tx *transaction.CLOLCLocalPlain) error {
	// Submit the transaction to the local chain
	sha256Func := sha256.New()
	sha256Func.Write([]byte(tx.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, err := crypto.PedersenCommit(tx.Amount)
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
	if err = o.SubmitTXLocalChain(clolcHidden); err != nil {
		return err
	}
	// Accumulate the commitment to the corresponding accumulator
	accumulatorKey := ComposeOrgRandMapKey(o.ID, TypeID(tx.CounterParty))
	if _, ok := o.epochAccumulatorMap[accumulatorKey]; !ok {
		o.epochAccumulatorMap[accumulatorKey] = commitment
	} else {
		o.epochAccumulatorMap[accumulatorKey].Add(
			o.epochAccumulatorMap[accumulatorKey],
			commitment,
		)
	}
	return nil
}

func (o *Organization) SubmitTXLocalChain(tx *transaction.CLOLCLocalHidden) error {
	return nil
}

func ComposeOrgRandMapKey(orgID1, orgID2 TypeID) [2]TypeID {
	if orgID1 < orgID2 {
		return [2]TypeID{orgID1, orgID2}
	}
	return [2]TypeID{orgID2, orgID1}
}
