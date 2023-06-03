package committee

import (
	"errors"

	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/organization"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

type typeID string

type Committee struct {
	ID                typeID
	managedEntityMap  map[auditor.TypeID][]organization.TypeID
	managedAuditorIDs []auditor.TypeID
	managedOrgIDs     []organization.TypeID
	epochTXRandMap    map[[2]string][]kyber.Scalar
	epochSecretKeyMap map[string]crypto.TypePrivateKey
	epochPublicKeyMap map[string]crypto.TypePublicKey
	epochOrgIDMap     map[organization.TypeID][]byte
	epochAuditorIDMap map[auditor.TypeID][]byte
}

func New(id string, auditors []*auditor.Auditor) *Committee {
	com := &Committee{
		ID:               typeID(id),
		managedEntityMap: make(map[auditor.TypeID][]organization.TypeID),
	}
	com.resetMaps()
	com.managedAuditorIDs = make([]auditor.TypeID, len(auditors))
	for idx, aud := range auditors {
		com.managedEntityMap[aud.ID] = aud.AuditedOrgIDs
		com.managedAuditorIDs[idx] = aud.ID
		com.managedOrgIDs = append(com.managedOrgIDs, aud.AuditedOrgIDs...)

	}
	return com
}

func (c *Committee) resetMaps() {
	c.epochTXRandMap = make(map[[2]string][]kyber.Scalar)
	c.epochSecretKeyMap = make(map[string]crypto.TypePrivateKey)
	c.epochPublicKeyMap = make(map[string]crypto.TypePublicKey)
	c.epochOrgIDMap = make(map[organization.TypeID][]byte)
	c.epochAuditorIDMap = make(map[auditor.TypeID][]byte)
}

// InitializeEpoch initialize the parameters for an auditing epoch
func (c *Committee) InitializeEpoch(
	auditors []*auditor.Auditor, organizations []*organization.Organization,
) (map[string]crypto.TypePublicKey, error) {
	c.resetMaps()
	// IN.1: generate randomness for the transactions {r_{i, j, k}},
	// note that r_{i, j, k} = r_{j, i, k}, and R_{i, j} = {r_{i, j, k}}_k
	if err := c.generateEpochTXRandomness(); err != nil {
		return nil, err
	}

	// IN.2: generate epoch random IDs for the organizations {id_i}
	if err := c.generateEpochOrgRandIDs(); err != nil {
		return nil, err
	}
	// IN.2: generate epoch random IDs for the auditors {id_z}
	if err := c.generateEpochAuditorRandIDs(); err != nil {
		return nil, err
	}

	// IN.3: generate secret-public key pairs for the organizations
	if err := c.generateEpochKeyPairs(); err != nil {
		return nil, err
	}

	// IN.4: publish the public keys (we just return the list of public keys at the end)

	// IN.5: forward the transaction randomnesses,auditor randomnesses and secret keys to the auditors
	// We need to forward: {r_{i, j, k}}, {r_z}, and {sk_i}
	for _, aud := range auditors {
		err := c.ForwardEpochAuditorParameters(aud)
		if err != nil {
			return nil, err
		}
	}
	// IN.5 (the missing step): forward the epoch random IDs to the organizations
	for _, org := range organizations {
		err := c.ForwardEpochOrgParameters(org)
		if err != nil {
			return nil, err
		}
	}

	// IN.4
	return c.epochPublicKeyMap, nil
}

func (c *Committee) generateEpochTXRandomness() error {
	// generate randomness for the transactions
	// distribute randomness to organizations
	// the complexity is O(n^2) here
	c.epochTXRandMap = make(map[[2]string][]kyber.Scalar)
	for i := 0; i < len(c.managedOrgIDs); i++ {
		for j := i + 1; j < len(c.managedOrgIDs); j++ {
			orgID1 := c.managedOrgIDs[i]
			orgID2 := c.managedOrgIDs[j]
			orgIDHash1 := organization.IDHashString(orgID1)
			orgIDHash2 := organization.IDHashString(orgID2)
			key := organization.IDHashKey(orgIDHash1, orgIDHash2)
			if _, ok := c.epochTXRandMap[key]; ok {
				continue
			}
			c.epochTXRandMap[key] = crypto.RandScalars(constants.MaxNumTXInEpoch)
		}
	}
	return nil
}

func (c *Committee) generateEpochKeyPairs() error {
	for _, id := range c.managedOrgIDs {
		privateKey, publicKey, err := crypto.KeyGen()
		if err != nil {
			return err
		}
		idHash := organization.IDHashString(id)
		c.epochSecretKeyMap[idHash] = privateKey
		c.epochPublicKeyMap[idHash] = publicKey
	}
	return nil
}

func (c *Committee) PublishPublicKeys() map[string]kyber.Point {
	publicKeyMap := make(map[string]kyber.Point)
	for _, id := range c.managedOrgIDs {
		idHash := organization.IDHashString(id)
		publicKeyMap[idHash] = c.epochPublicKeyMap[idHash]
	}
	return publicKeyMap
}

func (c *Committee) ForwardEpochAuditorParameters(auditor *auditor.Auditor) error {
	auditedOrgIDList, ok := c.managedEntityMap[auditor.ID]
	if !ok {
		return errors.New(string("auditor not found, id: " + auditor.ID))
	}
	// forward transaction randomnesses
	auditedOrgIDHashList := make([]string, len(auditedOrgIDList))
	for i, orgID := range auditedOrgIDList {
		auditedOrgIDHashList[i] = organization.IDHashString(orgID)
	}
	managedOrgIDHashList := make([]string, len(c.managedOrgIDs))
	for i, orgID := range c.managedOrgIDs {
		managedOrgIDHashList[i] = organization.IDHashString(orgID)
	}
	orgTXRandMap := make(map[[2]string][]kyber.Scalar)
	for _, orgIDHash1 := range auditedOrgIDHashList {
		for _, orgIDHash2 := range managedOrgIDHashList {
			if orgIDHash1 == orgIDHash2 {
				continue
			}
			key := organization.IDHashKey(orgIDHash1, orgIDHash2)
			if _, ok := c.epochTXRandMap[key]; !ok {
				return errors.New("randomness not found, key: " + key[0] + key[1])
			}
			orgTXRandMap[key] = c.epochTXRandMap[key]
		}
	}
	// Forward secret key
	auditedOrgSecretKeyMap := make(map[string]crypto.TypePrivateKey)
	for _, orgID := range auditedOrgIDList {
		orgIDHash := organization.IDHashString(orgID)
		secretKey, ok := c.epochSecretKeyMap[orgIDHash]
		if !ok {
			return errors.New(string("secret key not found, id: " + orgID))
		}
		auditedOrgSecretKeyMap[orgIDHash] = secretKey
	}
	auditor.SetEpochTXRandomness(orgTXRandMap)
	auditor.SetEpochSecretKey(auditedOrgSecretKeyMap)

	// set the epoch ID
	epochID, ok := c.epochAuditorIDMap[auditor.ID]
	if !ok {
		return errors.New(string("epoch ID not found, id: " + auditor.ID))
	}
	auditor.SetEpochID(epochID)

	// forward organization epoch ID
	epochOrgIDMap := make(map[organization.TypeID][]byte)
	for _, orgID := range auditedOrgIDList {
		epochID, ok := c.epochOrgIDMap[orgID]
		if !ok {
			return errors.New(string("epoch ID not found, id: " + orgID))
		}
		epochOrgIDMap[orgID] = epochID
	}
	auditor.SetEpochOrgIDMap(epochOrgIDMap)
	return nil
}

func (c *Committee) generateEpochOrgRandIDs() error {
	c.epochOrgIDMap = make(map[organization.TypeID][]byte)
	for _, id := range c.managedOrgIDs {
		randBytes, err := crypto.RandBytes()
		if err != nil {
			return err
		}
		c.epochOrgIDMap[id] = randBytes
	}
	return nil
}

func (c *Committee) generateEpochAuditorRandIDs() error {
	c.epochAuditorIDMap = make(map[auditor.TypeID][]byte)
	for _, id := range c.managedAuditorIDs {
		randBytes, err := crypto.RandBytes()
		if err != nil {
			return err
		}
		c.epochAuditorIDMap[id] = randBytes
	}
	return nil
}

func (c *Committee) ForwardEpochOrgParameters(org *organization.Organization) error {
	if _, ok := c.epochOrgIDMap[org.ID]; !ok {
		return errors.New(string("organization not found, id: " + org.ID))
	}
	org.SetEpochID(c.epochOrgIDMap[org.ID])
	return nil
}

func (c *Committee) VerifyOrgAndAudResult(
	orgID organization.TypeID,
	audID auditor.TypeID,
	orgChainTX *transaction.CLOLCOrgOnChain,
	audChainTX *transaction.CLOLCAudOnChain,
) (bool, error) {
	panic("implement me")
}

func (c *Committee) VerifyAuditPairResult(
	orgID1 organization.TypeID,
	orgID2 organization.TypeID,
	audID1 auditor.TypeID,
	audID2 auditor.TypeID,
	audChainTX1 *transaction.CLOLCAudOnChain,
	audChainTX2 *transaction.CLOLCAudOnChain,
) (bool, error) {
	panic("implement me")
}
