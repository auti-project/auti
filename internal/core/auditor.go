package core

import (
	"math/big"

	"go.dedis.ch/kyber/v3"
)

type Auditor struct {
	ID              string
	auditedOrgIDs   []string
	epochRand       *big.Int
	epochOrgRandMap map[[2]string][]*big.Int
	epochID         *big.Int
	epochSecretKey  kyber.Scalar
}

func NewAuditor(id string, organizations []*Organization) *Auditor {
	auditor := &Auditor{
		ID: id,
	}
	auditor.auditedOrgIDs = make([]string, len(organizations))
	for idx, org := range organizations {
		auditor.auditedOrgIDs[idx] = org.ID
	}
	return auditor
}

func (a *Auditor) SetEpochRandomness(random *big.Int) {
	a.epochRand = random
}

func (a *Auditor) SetEpochTXRandomness(txRandMap map[[2]string][]*big.Int) {
	a.epochOrgRandMap = txRandMap
}

func (a *Auditor) SetEpochSecretKey(secretKey kyber.Scalar) {
	a.epochSecretKey = secretKey
}
