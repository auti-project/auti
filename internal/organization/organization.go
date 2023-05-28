package organization

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

var (
	sha256Func = sha256.New()
)

type TypeID string

type Organization struct {
	ID                  TypeID
	IDHash              string
	epochRandID         *big.Int
	epochAccumulatorMap map[[2]string]kyber.Point
	epochTXRandomness   map[[2]string]kyber.Scalar
}

func New(id string) *Organization {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	idHash := hex.EncodeToString(sha256Func.Sum(nil))
	org := &Organization{
		ID:                  TypeID(id),
		IDHash:              idHash,
		epochAccumulatorMap: make(map[[2]string]kyber.Point),
	}
	return org
}

func (o *Organization) SetEpochID(randID *big.Int) {
	o.epochRandID = randID
}

func IDHashString(id TypeID) string {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	return hex.EncodeToString(sha256Func.Sum(nil))
}

func (o *Organization) RecordTransaction(tx *transaction.CLOLCLocalPlain) error {
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
	if err = o.SubmitTXLocalChain(clolcHidden); err != nil {
		return err
	}
	counterPartyHashStr := hex.EncodeToString(counterPartyHash)
	orgMapKey := IDHashKey(o.IDHash, counterPartyHashStr)
	// Accumulate the commitment to the corresponding accumulator
	if _, ok := o.epochAccumulatorMap[orgMapKey]; !ok {
		o.epochAccumulatorMap[orgMapKey] = commitment
	} else {
		o.epochAccumulatorMap[orgMapKey].Add(
			o.epochAccumulatorMap[orgMapKey],
			commitment,
		)
	}
	// Record the randomness used in the commitment
	o.epochTXRandomness[orgMapKey] = randScalar
	return nil
}

func (o *Organization) SubmitTXLocalChain(tx *transaction.CLOLCLocalHidden) error {
	return nil
}

//func IDKey(orgID1, orgID2 TypeID) [2]TypeID {
//	if orgID1 < orgID2 {
//		return [2]TypeID{orgID1, orgID2}
//	}
//	return [2]TypeID{orgID2, orgID1}
//}

func IDHashKey(orgIDHash1, orgIDHash2 string) [2]string {
	if orgIDHash1 < orgIDHash2 {
		return [2]string{orgIDHash1, orgIDHash2}
	}
	return [2]string{orgIDHash2, orgIDHash1}
}
