package closc

import (
	"go.dedis.ch/kyber/v3"

	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/organization"
	closcorg "github.com/auti-project/auti/internal/organization/closc"
	"github.com/auti-project/auti/internal/transaction/closc"
)

type Auditor struct {
	ID            auditor.TypeID
	AuditedOrgIDs []organization.TypeID
	EpochID       auditor.TypeEpochID
}

func New(id string, organizations []*closcorg.Organization) *Auditor {
	aud := &Auditor{
		ID: auditor.TypeID(id),
	}
	aud.AuditedOrgIDs = make([]organization.TypeID, len(organizations))
	for idx, org := range organizations {
		aud.AuditedOrgIDs[idx] = org.ID
	}
	return aud
}

func (a *Auditor) SetEpochID(id []byte) {
	a.EpochID = id
}

func (a *Auditor) VerifyMerkleProof(tx closc.LocalOnChain) (uint, error) {
	txPlain, err := tx.ToPlain()
	if err != nil {
		return 0, err
	}
	merkleProof, err := crypto.MerkleProofUnmarshal(txPlain.MerkleProof)
	if err != nil {
		return 0, err
	}
	ok, err := crypto.VerifyMerkleProof(txPlain, merkleProof, txPlain.MerkleRoot)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return 1, nil
}

func (a *Auditor) SummarizeMerkleProofVerificationResults(verificationResults []uint) bool {
	if len(verificationResults) == 0 {
		return false
	}
	result := verificationResults[0]
	for i := 1; i < len(verificationResults); i++ {
		result *= verificationResults[i]
		if result == 0 {
			return false
		}
	}
	return true
}

func (a *Auditor) VerifyCommitment(commitmentList1, commitmentList2 [][]byte,
	hashPoints1, hashPoints2 []kyber.Point) (bool, error) {
	if len(commitmentList1) != len(commitmentList2) {
		return false, nil
	}
	if len(commitmentList1) != len(hashPoints1) {
		return false, nil
	}
	if len(commitmentList2) != len(hashPoints2) {
		return false, nil
	}
	sum := crypto.KyberSuite.Point().Null()
	for i := 0; i < len(commitmentList1); i++ {
		commitPoint := crypto.KyberSuite.Point()
		if err := commitPoint.UnmarshalBinary(commitmentList1[i]); err != nil {
			return false, err
		}
		sum = sum.Add(sum, commitPoint)
	}
	for i := 0; i < len(commitmentList2); i++ {
		commitPoint := crypto.KyberSuite.Point()
		if err := commitPoint.UnmarshalBinary(commitmentList2[i]); err != nil {
			return false, err
		}
		sum = sum.Add(sum, commitPoint)
	}
	for i := 0; i < len(hashPoints1); i++ {
		sum = sum.Sub(sum, hashPoints1[i])
	}
	for i := 0; i < len(hashPoints2); i++ {
		sum = sum.Sub(sum, hashPoints2[i])
	}
	// TODO: test this
	return sum.Equal(crypto.KyberSuite.Point().Null()), nil
}
