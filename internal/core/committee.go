package core

import (
	"errors"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

type Committee struct {
	ID                  string
	managedEntityMap    map[string][]string
	managedAuditorIDs   []string
	managedOrgIDs       []string
	epochTXRandMap      map[[2]string][]*big.Int
	epochAuditorRandMap map[string]*big.Int
	epochKeyPairMap     map[string]*KeyPair
	epochOrgIDMap       map[string]*big.Int
	epochAuditorIDMap   map[string]*big.Int
}

func NewCommittee(id string, auditors []*Auditor) *Committee {
	committee := &Committee{
		ID:               id,
		managedEntityMap: make(map[string][]string),
	}
	committee.managedAuditorIDs = make([]string, len(auditors))
	for idx, auditor := range auditors {
		committee.managedEntityMap[auditor.ID] = auditor.auditedOrgIDs
		committee.managedAuditorIDs[idx] = auditor.ID
		for _, org := range auditor.auditedOrgIDs {
			committee.managedOrgIDs = append(committee.managedOrgIDs, org)
		}
	}
	return committee
}

func (c *Committee) InitializeEpoch(auditors []*Auditor, organizations []*Organization) error {
	// randomness generation
	// secret-key, public-key generation
	// publish public-key
	// distribute randomness, secret-key to auditors
	// generate ID for each organization
	// generate ID for each auditor
	if err := c.generateEpochTXRandomness(); err != nil {
		return err
	}
	if err := c.generateEpochAuditorRandomness(); err != nil {
		return err
	}
	if err := c.generateEpochKeyPairs(); err != nil {
		return err
	}
	for _, auditor := range auditors {
		err := c.ForwardEpochParameters(auditor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Committee) generateEpochTXRandomness() error {
	// generate randomness for the transactions
	// distribute randomness to organizations
	// the complexity is O(n^2) here
	c.epochTXRandMap = make(map[[2]string][]*big.Int)
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

func composeOrgRandMapKey(orgID1, orgID2 string) [2]string {
	if orgID1 < orgID2 {
		return [2]string{orgID1, orgID2}
	}
	return [2]string{orgID2, orgID1}
}

func (c *Committee) generateEpochAuditorRandomness() error {
	// generate randomness for auditors
	// distribute randomness to auditors
	c.epochAuditorIDMap = make(map[string]*big.Int)
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
	c.epochKeyPairMap = make(map[string]*KeyPair)
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

func (c *Committee) PublishPublicKeys() map[string]kyber.Point {
	publicKeyMap := make(map[string]kyber.Point)
	for _, id := range c.managedOrgIDs {
		publicKeyMap[id] = c.epochKeyPairMap[id].PublicKey
	}
	return publicKeyMap
}

func (c *Committee) ForwardEpochParameters(auditor *Auditor) error {
	// forward transaction randomnesses
	auditedOrgIDList, ok := c.managedEntityMap[auditor.ID]
	if !ok {
		return errors.New("auditor not found, id: " + auditor.ID)
	}
	orgTXRandMap := make(map[[2]string][]*big.Int)
	for _, orgID1 := range auditedOrgIDList {
		for _, orgID2 := range c.managedOrgIDs {
			if orgID1 == orgID2 {
				continue
			}
			key := composeOrgRandMapKey(orgID1, orgID2)
			if _, ok := c.epochTXRandMap[key]; !ok {
				return errors.New("randomness not found, key: " + key[0] + key[1])
			}
			orgTXRandMap[key] = c.epochTXRandMap[key]
		}
	}
	auditor.SetEpochTXRandomness(orgTXRandMap)
	// forward auditor randomness
	epochAuditorRand, ok := c.epochAuditorRandMap[auditor.ID]
	if !ok {
		return errors.New("auditor randomness not found, id: " + auditor.ID)
	}
	auditor.SetEpochRandomness(epochAuditorRand)
	// forward secret key
	keyPair, ok := c.epochKeyPairMap[auditor.ID]
	if !ok {
		return errors.New("secret key not found, id: " + auditor.ID)
	}
	auditor.SetEpochSecretKey(keyPair.PrivateKey)
	return nil
}

func (c *Committee) generateEpochOrgRandIDs() error {
	c.epochOrgIDMap = make(map[string]*big.Int)
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
	c.epochAuditorIDMap = make(map[string]*big.Int)
	for _, id := range c.managedAuditorIDs {
		randInt, err := RandInt()
		if err != nil {
			return err
		}
		c.epochAuditorIDMap[id] = randInt
	}
	return nil
}
