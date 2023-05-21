package auditor

import (
	"math/big"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/organization"
)

type TypeID string

type Auditor struct {
	ID                   TypeID
	AuditedOrgIDs        []organization.TypeID
	epochRand            *big.Int
	epochOrgRandMap      map[[2]organization.TypeID][]*big.Int
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

func (a *Auditor) SetEpochTXRandomness(txRandMap map[[2]organization.TypeID][]*big.Int) {
	a.epochOrgRandMap = txRandMap
}

func (a *Auditor) SetEpochSecretKey(orgSecretKeyMap map[organization.TypeID]crypto.TypeSecretKey) {
	a.epochOrgSecretKeyMap = orgSecretKeyMap
}

func (a *Auditor) SetEpochID(id *big.Int) {
	a.epochID = id
}
