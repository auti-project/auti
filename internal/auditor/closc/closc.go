package closc

import (
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/organization"
	closcorg "github.com/auti-project/auti/internal/organization/closc"
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
