package auditor

import (
	"crypto/sha256"
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
	epochOrgSecretKeyMap map[string]crypto.TypePrivateKey
	epochOrgIDMap        map[organization.TypeID]*big.Int
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

func (a *Auditor) SetEpochSecretKey(orgSecretKeyMap map[string]crypto.TypePrivateKey) {
	a.epochOrgSecretKeyMap = orgSecretKeyMap
}

func (a *Auditor) SetEpochID(id *big.Int) {
	a.epochID = id
}

func (a *Auditor) SetEpochOrgIDMap(idMap map[organization.TypeID]*big.Int) {
	a.epochOrgIDMap = idMap
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

func (a *Auditor) ComputeB(orgTXRandList, comTXRandList []kyber.Scalar) (kyber.Point, error) {
	if len(orgTXRandList) != len(comTXRandList) {
		return nil, fmt.Errorf("length of two lists are not equal")
	}
	scalar := crypto.KyberSuite.Scalar().Zero()
	for idx := range orgTXRandList {
		tmp := crypto.KyberSuite.Scalar().Mul(orgTXRandList[idx], comTXRandList[idx])
		scalar.Sub(scalar, tmp)
	}
	result := crypto.KyberSuite.Point().Mul(scalar, crypto.KyberSuite.Point().Base())
	return result, nil
}

func (a *Auditor) ComputeC(res, A kyber.Point) kyber.Point {
	result := crypto.KyberSuite.Point().Add(res, A)
	return result
}
func (a *Auditor) ComputeD(pointA, pointB kyber.Point) kyber.Point {
	negA := crypto.KyberSuite.Point().Neg(pointA)
	negB := crypto.KyberSuite.Point().Neg(pointB)
	result := crypto.KyberSuite.Point().Add(negA, negB)
	return result
}

func (a *Auditor) EncryptConsistencyExamResult(orgID organization.TypeID, counterPartyIDHash string,
	res, pointB, pointC, pointD kyber.Point, publicKey kyber.Point) (*transaction.CLOLCAudPlain, error) {
	txID, err := a.ComputeCETransactionID(orgID, counterPartyIDHash)
	if err != nil {
		return nil, err
	}
	cipherRes, err := crypto.EncryptPoint(publicKey, res)
	if err != nil {
		return nil, err
	}
	cipherResBytes, err := cipherRes.Serialize()
	if err != nil {
		return nil, err
	}
	cipherB, err := crypto.EncryptPoint(publicKey, pointB)
	if err != nil {
		return nil, err
	}
	cipherBBytes, err := cipherB.Serialize()
	if err != nil {
		return nil, err
	}
	cipherC, err := crypto.EncryptPoint(publicKey, pointC)
	if err != nil {
		return nil, err
	}
	cipherCBytes, err := cipherC.Serialize()
	if err != nil {
		return nil, err
	}
	cipherD, err := crypto.EncryptPoint(publicKey, pointD)
	if err != nil {
		return nil, err
	}
	cipherDBytes, err := cipherD.Serialize()
	if err != nil {
		return nil, err
	}
	return transaction.NewCLOLCAudPlain(txID, cipherResBytes, cipherBBytes, cipherCBytes, cipherDBytes), nil
}

func (a *Auditor) ComputeCETransactionID(orgID organization.TypeID, counterPartyIDHash string) ([]byte, error) {
	orgIDHashStr := organization.IDHashString(orgID)
	orgKey := organization.IDHashKey(orgIDHashStr, counterPartyIDHash)
	randomnesses := a.epochOrgRandMap[orgKey]
	epochOrgID := a.epochOrgIDMap[orgID]
	epochOrgIDBytes := epochOrgID.Bytes()
	randAccumulator := crypto.KyberSuite.Scalar().Zero()
	for _, randScalar := range randomnesses {
		randAccumulator.Add(randAccumulator, randScalar)
	}
	randAccumulatorBytes, err := randAccumulator.MarshalBinary()
	if err != nil {
		return nil, err
	}
	concatBytes := append(epochOrgIDBytes, randAccumulatorBytes...)
	sha256Func := sha256.New()
	sha256Func.Write(concatBytes)
	result := sha256Func.Sum(nil)
	return result, nil
}
