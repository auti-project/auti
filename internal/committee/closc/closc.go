package closc

import (
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/auditor/clolc"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/organization"
)

type Committee struct {
	ID                committee.TypeID
	managedEntityMap  map[auditor.TypeID][]organization.TypeID
	managedAuditorIDs []auditor.TypeID
	managedOrgIDs     []organization.TypeID
	epochAuditorIDMap map[auditor.TypeID]auditor.TypeEpochID
}

func New(id string, auditors []*clolc.Auditor) *Committee {
	com := &Committee{
		ID:               committee.TypeID(id),
		managedEntityMap: make(map[auditor.TypeID][]organization.TypeID),
	}
	com.managedAuditorIDs = make([]auditor.TypeID, len(auditors))
	for idx, aud := range auditors {
		com.managedEntityMap[aud.ID] = aud.AuditedOrgIDs
		com.managedAuditorIDs[idx] = aud.ID
		com.managedOrgIDs = append(com.managedOrgIDs, aud.AuditedOrgIDs...)
	}
	return com
}

func (c *Committee) reinitializeMaps() {
	c.epochAuditorIDMap = make(map[auditor.TypeID]auditor.TypeEpochID)
}

//func (c *CLOSCCommittee) InitializeEpoch(auditors []*auditor.CLOSC)
