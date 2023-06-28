package closc

import (
	"crypto/rand"

	"github.com/auti-project/auti/internal/auditor"
	closcaud "github.com/auti-project/auti/internal/auditor/closc"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/organization"
)

type Committee struct {
	ID                committee.TypeID
	managedEntityMap  map[auditor.TypeID][]organization.TypeID
	managedAuditorIDs []auditor.TypeID
	managedOrgIDs     []organization.TypeID
	epochAuditorIDMap map[auditor.TypeID]auditor.TypeEpochID
}

func New(id string, auditors []*closcaud.Auditor) *Committee {
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

func (c *Committee) InitializeEpoch(auditors []*closcaud.Auditor) error {
	c.reinitializeMaps()
	for _, aud := range auditors {
		// Generate epoch ID for each auditor
		epochIDBytes := make([]byte, constants.SecurityParameterBytes)
		_, err := rand.Read(epochIDBytes)
		if err != nil {
			return err
		}
		c.epochAuditorIDMap[aud.ID] = epochIDBytes
		// Distribute epoch auditor IDs
		aud.SetEpochID(epochIDBytes)
	}
	return nil
}
