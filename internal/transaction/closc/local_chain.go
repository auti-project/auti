package closc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type LocalCommitmentPlain struct {
	Commitment []byte
}

func NewLocalCommitmentPlain(commitment []byte) *LocalCommitmentPlain {
	return &LocalCommitmentPlain{
		Commitment: commitment,
	}
}

func (l *LocalCommitmentPlain) ToOnChain() *LocalCommitmentOnChain {
	return NewLocalCommitmentOnChain(hex.EncodeToString(l.Commitment))
}

type LocalCommitmentOnChain struct {
	Commitment string `json:"commitment"`
}

func NewLocalCommitmentOnChain(commitment string) *LocalCommitmentOnChain {
	return &LocalCommitmentOnChain{
		Commitment: commitment,
	}
}

func (l *LocalCommitmentOnChain) ToPlain() (*LocalCommitmentPlain, error) {
	commitment, err := hex.DecodeString(l.Commitment)
	if err != nil {
		return nil, err
	}
	return NewLocalCommitmentPlain(commitment), nil
}

func (l *LocalCommitmentOnChain) KeyVal() (string, []byte, error) {
	txJSON, err := json.Marshal(l)
	if err != nil {
		return "", nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(txJSON)
	return hex.EncodeToString(sha256Func.Sum(nil)), txJSON, nil
}

type LocalPlain struct {
	Commitment  []byte
	MerkleRoot  []byte
	MerkleProof []byte
}

// TODO: work on this after the merkle tree is implemented
