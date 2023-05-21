package core

import (
	"errors"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

type CommitteeID string

type Committee struct {
	ID                  CommitteeID
	managedEntityMap    map[AuditorID][]OrgID
	managedAuditorIDs   []AuditorID
	managedOrgIDs       []OrgID
	epochTXRandMap      map[[2]OrgID][]*big.Int
	epochAuditorRandMap map[AuditorID]*big.Int
	epochKeyPairMap     map[OrgID]*KeyPair
	epochOrgIDMap       map[OrgID]*big.Int
	epochAuditorIDMap   map[AuditorID]*big.Int
}

func NewCommittee(id string, auditors []*Auditor) *Committee {
	committee := &Committee{
		ID:               CommitteeID(id),
		managedEntityMap: make(map[AuditorID][]OrgID),
	}
	committee.managedAuditorIDs = make([]AuditorID, len(auditors))
	for idx, auditor := range auditors {
		committee.managedEntityMap[auditor.ID] = auditor.auditedOrgIDs
		committee.managedAuditorIDs[idx] = auditor.ID
		for _, org := range auditor.auditedOrgIDs {
			committee.managedOrgIDs = append(committee.managedOrgIDs, org)
		}
	}
	return committee
}

// InitializeEpoch initialize the parameters for an auditing epoch
func (c *Committee) InitializeEpoch(
	auditors []*Auditor, organizations []*Organization,
) (map[OrgID]TypePublicKey, error) {
	// randomness generation
	// secret-key, public-key generation
	// publish public-key
	// distribute randomness, secret-key to auditors
	// generate ID for each organization
	// generate ID for each auditor

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

	// IN.5: forward the transaction randomnesses, auditor randomnesses and secret keys to the auditors
	// We need to forward: {r_{i, j, k}}, {r_z}, and {sk_i}
	for _, auditor := range auditors {
		err := c.ForwardEpochParameters(auditor)
		if err != nil {
			return nil, err
		}
	}

	orgPublicKeyMap := make(map[OrgID]TypePublicKey)
	for _, org := range organizations {
		orgPublicKeyMap[org.ID] = c.epochKeyPairMap[org.ID].PublicKey
	}
	return orgPublicKeyMap, nil
}

func (c *Committee) generateEpochTXRandomness() error {
	// generate randomness for the transactions
	// distribute randomness to organizations
	// the complexity is O(n^2) here
	c.epochTXRandMap = make(map[[2]OrgID][]*big.Int)
	for i := 0; i < len(c.managedOrgIDs); i++ {
		for j := i + 1; j < len(c.managedOrgIDs); j++ {
			orgID1 := c.managedOrgIDs[i]
			orgID2 := c.managedOrgIDs[j]
			key := composeOrgRandMapKey(orgID1, orgID2)
			if _, ok := c.epochTXRandMap[key]; ok {
				continue
			}
			randList, err := RandIntList(MaxNumTXInEpoch)
			if err != nil {
				return err
			}
			c.epochTXRandMap[key] = randList
		}
	}
	return nil
}

func composeOrgRandMapKey(orgID1, orgID2 OrgID) [2]OrgID {
	if orgID1 < orgID2 {
		return [2]OrgID{orgID1, orgID2}
	}
	return [2]OrgID{orgID2, orgID1}
}

func (c *Committee) generateEpochAuditorRandomness() error {
	// generate randomness for auditors
	// distribute randomness to auditors
	c.epochAuditorIDMap = make(map[AuditorID]*big.Int)
	for _, id := range c.managedAuditorIDs {
		randInt, err := RandInt()
		if err != nil {
			return err
		}
		c.epochAuditorRandMap[id] = randInt
	}
	return nil
}

func (c *Committee) generateEpochKeyPairs() error {
	c.epochKeyPairMap = make(map[OrgID]*KeyPair)
	for _, id := range c.managedOrgIDs {
		privateKey, publicKey, err := KeyGen()
		if err != nil {
			return err
		}
		c.epochKeyPairMap[id] = &KeyPair{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}
	}
	return nil
}

func (c *Committee) PublishPublicKeys() map[OrgID]kyber.Point {
	publicKeyMap := make(map[OrgID]kyber.Point)
	for _, id := range c.managedOrgIDs {
		publicKeyMap[id] = c.epochKeyPairMap[id].PublicKey
	}
	return publicKeyMap
}

func (c *Committee) ForwardEpochParameters(auditor *Auditor) error {
	auditedOrgIDList, ok := c.managedEntityMap[auditor.ID]
	if !ok {
		return errors.New(string("auditor not found, id: " + auditor.ID))
	}
	// forward transaction randomnesses, and at the meantime, forward the secret keys
	orgTXRandMap := make(map[[2]OrgID][]*big.Int)
	auditedOrgSecretKeyMap := make(map[OrgID]TypeSecretKey)
	for _, orgID1 := range auditedOrgIDList {
		// check if both the organizations in the pair are audited by the auditor
		for _, orgID2 := range c.managedOrgIDs {
			if orgID1 == orgID2 {
				continue
			}
			key := composeOrgRandMapKey(orgID1, orgID2)
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
	return nil
}

func (c *Committee) generateEpochOrgRandIDs() error {
	c.epochOrgIDMap = make(map[OrgID]*big.Int)
	for _, id := range c.managedOrgIDs {
		randInt, err := RandInt()
		if err != nil {
			return err
		}
		c.epochOrgIDMap[id] = randInt
	}
	return nil
}

func (c *Committee) generateEpochAuditorRandIDs() error {
	c.epochAuditorIDMap = make(map[AuditorID]*big.Int)
	for _, id := range c.managedAuditorIDs {
		randInt, err := RandInt()
		if err != nil {
			return err
		}
		c.epochAuditorIDMap[id] = randInt
	}
	return nil
}
