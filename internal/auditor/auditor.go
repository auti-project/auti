package auditor

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/organization"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

type TypeID string

type Auditor struct {
	ID                   TypeID
	AuditedOrgIDs        []organization.TypeID
	epochRand            *big.Int
	epochOrgRandMap      map[[2]string][]kyber.Scalar
	epochID              *big.Int
	epochOrgSecretKeyMap map[organization.TypeID]crypto.TypeSecretKey
}

func New(id string, organizations []*organization.Organization) *Auditor {
	auditor := &Auditor{
		ID: TypeID(id),
	}
	auditor.AuditedOrgIDs = make([]organization.TypeID, len(organizations))
	for idx, org := range organizations {
		auditor.AuditedOrgIDs[idx] = org.ID
	}
	return auditor
}

func (a *Auditor) SetEpochRandomness(random *big.Int) {
	a.epochRand = random
}

func (a *Auditor) SetEpochTXRandomness(txRandMap map[[2]string][]kyber.Scalar) {
	a.epochOrgRandMap = txRandMap
}

func (a *Auditor) SetEpochSecretKey(orgSecretKeyMap map[organization.TypeID]crypto.TypeSecretKey) {
	a.epochOrgSecretKeyMap = orgSecretKeyMap
}

func (a *Auditor) SetEpochID(id *big.Int) {
	a.epochID = id
}

func (a *Auditor) AccumulateCommitments(
	orgID organization.TypeID, txList []*transaction.CLOLCLocalHidden,
) (kyber.Point, error) {
	if len(txList) == 0 {
		return nil, fmt.Errorf("empty transaction list")
	}
	if constants.MaxNumTXInEpoch < uint(len(txList)) {
		return nil, fmt.Errorf("too many transactions in the epoch: %d", len(txList))
	}
	orgIDHashStr := organization.IDHashString(orgID)
	counterPartyIDHashStr := hex.EncodeToString(txList[0].CounterParty)
	orgKey := organization.IDHashKey(orgIDHashStr, counterPartyIDHashStr)
	randomScalars := a.epochOrgRandMap[orgKey]
	result := crypto.KyberSuite.Point().Null()
	for idx, tx := range txList {
		commitmentBytes := tx.Commitment
		commitmentPoint := crypto.KyberSuite.Point()
		if err := commitmentPoint.UnmarshalBinary(commitmentBytes); err != nil {
			return nil, err
		}
		commitmentPoint.Mul(randomScalars[idx], commitmentPoint)
		result.Add(result, commitmentPoint)
	}
	return result, nil
}

func (a *Auditor) ComputeB(orgTXRandList, comTXRandList []kyber.Scalar) (kyber.Scalar, error) {
	if len(orgTXRandList) != len(comTXRandList) {
		return nil, fmt.Errorf("length of two lists are not equal")
	}
	result := crypto.KyberSuite.Scalar().Zero()
	for idx := range orgTXRandList {
		tmp := crypto.KyberSuite.Scalar().Mul(orgTXRandList[idx], comTXRandList[idx])
		result.Sub(result, tmp)
	}
	return result, nil
}
