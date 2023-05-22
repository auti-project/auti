package organization

import (
	"crypto/sha256"
	"math/big"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/transaction"
)

type TypeID string

type Organization struct {
	ID                  TypeID
	epochRandID         *big.Int
	epochAccumulatorMap map[[2]TypeID]*big.Int
}

func New(id string) *Organization {
	org := &Organization{
		ID:                  TypeID(id),
		epochAccumulatorMap: make(map[[2]TypeID]*big.Int),
	}
	return org
}

func (o *Organization) SetEpochID(randID *big.Int) {
	o.epochRandID = randID
}

func (o *Organization) RecordTransaction(tx *transaction.CLOLCPlain) error {
	// Submit the transaction to the local chain
	sha256Func := sha256.New()
	sha256Func.Write([]byte(tx.CounterParty))
	counterPartyHash := sha256Func.Sum(nil)
	commitment, err := crypto.PedersenCommit(tx.Amount)
	if err != nil {
		return err
	}
	clolcCipher := &transaction.CLOLCCipher{
		CounterParty: counterPartyHash,
		Commitment:   commitment,
	}
	if err = o.SubmitTXLocalChain(clolcCipher); err != nil {
		return err
	}
	// Accumulate the commitment to the corresponding accumulator
	accumulatorKey := ComposeOrgRandMapKey(o.ID, TypeID(tx.CounterParty))
	if _, ok := o.epochAccumulatorMap[accumulatorKey]; !ok {
		o.epochAccumulatorMap[accumulatorKey] = big.NewInt(1)
	}
	o.epochAccumulatorMap[accumulatorKey].Mul(
		o.epochAccumulatorMap[accumulatorKey],
		new(big.Int).SetBytes(commitment),
	)
	return nil
}

func (o *Organization) SubmitTXLocalChain(tx *transaction.CLOLCCipher) error {
	return nil
}

func ComposeOrgRandMapKey(orgID1, orgID2 TypeID) [2]TypeID {
	if orgID1 < orgID2 {
		return [2]TypeID{orgID1, orgID2}
	}
	return [2]TypeID{orgID2, orgID1}
}
