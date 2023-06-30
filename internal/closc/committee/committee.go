package committee

import (
	"crypto/rand"

	"github.com/auti-project/auti/internal/closc/auditor"
	closcorg "github.com/auti-project/auti/internal/closc/organization"
	"github.com/auti-project/auti/internal/constants"
)

type TypeID string

type Committee struct {
	ID                TypeID
	managedEntityMap  map[auditor.TypeID][]closcorg.TypeID
	managedAuditorIDs []auditor.TypeID
	managedOrgIDs     []closcorg.TypeID
	epochAuditorIDMap map[auditor.TypeID]auditor.TypeEpochID
}

func New(id string, auditors []*auditor.Auditor) *Committee {
	com := &Committee{
		ID:               TypeID(id),
		managedEntityMap: make(map[auditor.TypeID][]closcorg.TypeID),
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

func (c *Committee) InitializeEpoch(auditors []*auditor.Auditor) error {
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

// TODO: Implement this
//func (c *Committee) VerifyMerkleBatchProof()