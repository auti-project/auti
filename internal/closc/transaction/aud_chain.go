package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"go.dedis.ch/kyber/v3"

	"github.com/auti-project/auti/internal/crypto"
)

type AudPlain struct {
	Commitment []byte
	BatchProof []byte
}

func NewAudPlain(commitment []byte, batchProof []byte) *AudPlain {
	return &AudPlain{
		Commitment: commitment,
		BatchProof: batchProof,
	}
}

func NewAudPlainFromPointAndProof(commitment kyber.Point, batchProof *crypto.MerkleBatchProof) (*AudPlain, error) {
	commitmentBytes, err := commitment.MarshalBinary()
	if err != nil {
		return nil, err
	}
	batchProofBytes, err := crypto.MerkleBatchProofMarshal(batchProof)
	if err != nil {
		return nil, err
	}
	return NewAudPlain(commitmentBytes, batchProofBytes), nil
}

func (a *AudPlain) ToOnChain() *AudOnChain {
	return NewAudOnChain(
		hex.EncodeToString(a.Commitment),
		hex.EncodeToString(a.BatchProof),
	)
}

type AudOnChain struct {
	Commitment string `json:"commitment"`
	BatchProof string `json:"batch_proof"`
}

func NewAudOnChain(commitment string, batchProof string) *AudOnChain {
	return &AudOnChain{
		Commitment: commitment,
		BatchProof: batchProof,
	}
}

func (a *AudOnChain) ToPlain() (*AudPlain, error) {
	commitment, err := hex.DecodeString(a.Commitment)
	if err != nil {
		return nil, err
	}
	batchProof, err := hex.DecodeString(a.BatchProof)
	if err != nil {
		return nil, err
	}
	return NewAudPlain(commitment, batchProof), nil
}

func (a *AudOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(a)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}
