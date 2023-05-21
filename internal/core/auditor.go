package core

import (
	"math/big"
)

type AuditorID string

type Auditor struct {
	ID                   AuditorID
	auditedOrgIDs        []OrgID
	epochRand            *big.Int
	epochOrgRandMap      map[[2]OrgID][]*big.Int
	epochID              *big.Int
	epochOrgSecretKeyMap map[OrgID]TypeSecretKey
}

func NewAuditor(id string, organizations []*Organization) *Auditor {
	auditor := &Auditor{
		ID: AuditorID(id),
	}
	auditor.auditedOrgIDs = make([]OrgID, len(organizations))
	for idx, org := range organizations {
		auditor.auditedOrgIDs[idx] = org.ID
	}
	return auditor
}

func (a *Auditor) SetEpochRandomness(random *big.Int) {
	a.epochRand = random
}

func (a *Auditor) SetEpochTXRandomness(txRandMap map[[2]OrgID][]*big.Int) {
	a.epochOrgRandMap = txRandMap
}

func (a *Auditor) SetEpochSecretKey(orgSecretKeyMap map[OrgID]TypeSecretKey) {
	a.epochOrgSecretKeyMap = orgSecretKeyMap
}
