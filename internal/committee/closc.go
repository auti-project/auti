package committee

import (
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/organization"
)

type CLOSCCommittee struct {
	ID                typeID
	managedEntityMap  map[auditor.TypeID][]organization.TypeID
	managedAuditorIDs []auditor.TypeID
	managedOrgIDs     []organization.TypeID
	epochAuditorIDMap map[auditor.TypeID]auditor.TypeEpochID
}

func CLOSCNew(id string, auditors []*auditor.CLOLCAuditor) *CLOSCCommittee {
	com := &CLOSCCommittee{
		ID:               typeID(id),
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

func (c *CLOSCCommittee) reinitializeMaps() {
	c.epochAuditorIDMap = make(map[auditor.TypeID]auditor.TypeEpochID)
}

//func (c *CLOSCCommittee) InitializeEpoch(auditors []*auditor.CLOSC)
