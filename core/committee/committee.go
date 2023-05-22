package committee

import (
	"errors"
	"math/big"

	"github.com/auti-project/auti/core/auditor"
	"github.com/auti-project/auti/core/constants"
	"github.com/auti-project/auti/core/crypto"
	"github.com/auti-project/auti/core/organization"
	"go.dedis.ch/kyber/v3"
)

type typeID string

type Committee struct {
	ID                  typeID
	managedEntityMap    map[auditor.TypeID][]organization.TypeID
	managedAuditorIDs   []auditor.TypeID
	managedOrgIDs       []organization.TypeID
	epochTXRandMap      map[[2]organization.TypeID][]*big.Int
	epochAuditorRandMap map[auditor.TypeID]*big.Int
	epochKeyPairMap     map[organization.TypeID]*crypto.KeyPair
	epochOrgIDMap       map[organization.TypeID]*big.Int
	epochAuditorIDMap   map[auditor.TypeID]*big.Int
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
		for _, org := range aud.AuditedOrgIDs {
			com.managedOrgIDs = append(com.managedOrgIDs, org)
		}
	}
	return com
}

func (c *Committee) resetMaps() {
	c.epochTXRandMap = make(map[[2]organization.TypeID][]*big.Int)
	c.epochAuditorRandMap = make(map[auditor.TypeID]*big.Int)
	c.epochKeyPairMap = make(map[organization.TypeID]*crypto.KeyPair)
	c.epochOrgIDMap = make(map[organization.TypeID]*big.Int)
	c.epochAuditorIDMap = make(map[auditor.TypeID]*big.Int)
}

// InitializeEpoch initialize the parameters for an auditing epoch
func (c *Committee) InitializeEpoch(
	auditors []*auditor.Auditor, organizations []*organization.Organization,
) (map[organization.TypeID]crypto.TypePublicKey, error) {
	c.resetMaps()
	// IN.1: generate randomness for the transactions {r_{i, j, k}},
	// note that r_{i, j, k} = r_{j, i, k}, and R_{i, j} = {r_{i, j, k}}_k
	if err := c.generateEpochTXRandomness(); err != nil {
		return nil, err
	}
	// IN.1: generate randomness for the auditors {r_z}
	if err := c.generateEpochAuditorRandomness(); err != nil {
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
	orgPublicKeyMap := make(map[organization.TypeID]crypto.TypePublicKey)
	for _, orgID := range c.managedOrgIDs {
		orgPublicKeyMap[orgID] = c.epochKeyPairMap[orgID].PublicKey
	}
	return orgPublicKeyMap, nil
}

func (c *Committee) generateEpochTXRandomness() error {
	// generate randomness for the transactions
	// distribute randomness to organizations
	// the complexity is O(n^2) here
	c.epochTXRandMap = make(map[[2]organization.TypeID][]*big.Int)
	for i := 0; i < len(c.managedOrgIDs); i++ {
		for j := i + 1; j < len(c.managedOrgIDs); j++ {
			orgID1 := c.managedOrgIDs[i]
			orgID2 := c.managedOrgIDs[j]
			key := organization.ComposeOrgRandMapKey(orgID1, orgID2)
			if _, ok := c.epochTXRandMap[key]; ok {
				continue
			}
			randList, err := crypto.RandIntList(constants.MaxNumTXInEpoch)
			if err != nil {
				return err
			}
			c.epochTXRandMap[key] = randList
		}
	}
	return nil
}

func (c *Committee) generateEpochAuditorRandomness() error {
	// generate randomness for auditors
	// distribute randomness to auditors
	c.epochAuditorIDMap = make(map[auditor.TypeID]*big.Int)
	for _, id := range c.managedAuditorIDs {
		randInt, err := crypto.RandInt()
		if err != nil {
			return err
		}
		c.epochAuditorRandMap[id] = randInt
	}
	return nil
}

func (c *Committee) generateEpochKeyPairs() error {
	c.epochKeyPairMap = make(map[organization.TypeID]*crypto.KeyPair)
	for _, id := range c.managedOrgIDs {
		privateKey, publicKey, err := crypto.KeyGen()
		if err != nil {
			return err
		}
		c.epochKeyPairMap[id] = &crypto.KeyPair{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}
	}
	return nil
}

func (c *Committee) PublishPublicKeys() map[organization.TypeID]kyber.Point {
	publicKeyMap := make(map[organization.TypeID]kyber.Point)
	for _, id := range c.managedOrgIDs {
		publicKeyMap[id] = c.epochKeyPairMap[id].PublicKey
	}
	return publicKeyMap
}

func (c *Committee) ForwardEpochAuditorParameters(auditor *auditor.Auditor) error {
	auditedOrgIDList, ok := c.managedEntityMap[auditor.ID]
	if !ok {
		return errors.New(string("auditor not found, id: " + auditor.ID))
	}
	// forward transaction randomnesses, and at the meantime, forward the secret keys
	orgTXRandMap := make(map[[2]organization.TypeID][]*big.Int)
	auditedOrgSecretKeyMap := make(map[organization.TypeID]crypto.TypeSecretKey)
	for _, orgID1 := range auditedOrgIDList {
		for _, orgID2 := range c.managedOrgIDs {
			if orgID1 == orgID2 {
				continue
			}
			key := organization.ComposeOrgRandMapKey(orgID1, orgID2)
			if _, ok := c.epochTXRandMap[key]; !ok {
				return errors.New(string("randomness not found, key: " + key[0] + key[1]))
			}
			orgTXRandMap[key] = c.epochTXRandMap[key]
		}
		// forward secret key
		keyPair, ok := c.epochKeyPairMap[orgID1]
		if !ok {
			return errors.New(string("key pair not found, id: " + orgID1))
		}
		auditedOrgSecretKeyMap[orgID1] = keyPair.PrivateKey
	}
	auditor.SetEpochTXRandomness(orgTXRandMap)
	auditor.SetEpochSecretKey(auditedOrgSecretKeyMap)

	// forward auditor randomness
	epochAuditorRand, ok := c.epochAuditorRandMap[auditor.ID]
	if !ok {
		return errors.New(string("auditor randomness not found, id: " + auditor.ID))
	}
	auditor.SetEpochRandomness(epochAuditorRand)

	// set the epoch ID
	epochID, ok := c.epochAuditorIDMap[auditor.ID]
	if !ok {
		return errors.New(string("epoch ID not found, id: " + auditor.ID))
	}
	auditor.SetEpochID(epochID)
	return nil
}

func (c *Committee) generateEpochOrgRandIDs() error {
	c.epochOrgIDMap = make(map[organization.TypeID]*big.Int)
	for _, id := range c.managedOrgIDs {
		randInt, err := crypto.RandInt()
		if err != nil {
			return err
		}
		c.epochOrgIDMap[id] = randInt
	}
	return nil
}

func (c *Committee) generateEpochAuditorRandIDs() error {
	c.epochAuditorIDMap = make(map[auditor.TypeID]*big.Int)
	for _, id := range c.managedAuditorIDs {
		randInt, err := crypto.RandInt()
		if err != nil {
			return err
		}
		c.epochAuditorIDMap[id] = randInt
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
