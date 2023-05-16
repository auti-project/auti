package core

import "math/big"

type Committee struct {
	ID                   string
	ManagedAuditors      []*Auditor
	ManagedOrganizations []*Organization
	orgRandMap           map[[2]string][]*big.Int
	auditorRandList      []*big.Int
	PublicKeys           []*big.Int
}

func NewCommittee(auditors []*Auditor) *Committee {
	committee := &Committee{
		ManagedAuditors: auditors,
	}
	for _, auditor := range auditors {
		committee.ManagedOrganizations = append(committee.ManagedOrganizations, auditor.AuditedOrganizations...)
	}
	return committee
}

func (c *Committee) InitializeEpoch() error {
	// randomness generation
	// secret-key, public-key generation
	// publish public-key
	// distribute randomness, secret-key to auditors
	// generate ID for each organization
	// generate ID for each auditor
	if err := c.epochTransactionRandomness(); err != nil {
		return err
	}
	if err := c.epochAuditorRandomness(); err != nil {
		return err
	}
	for i := 0; i < len(c.ManagedAuditors); i++ {
	}
	return nil
}

func (c *Committee) epochTransactionRandomness() error {
	// generate randomness for the transactions
	// distribute randomness to organizations
	// the complexity is O(n^2) here
	for i := 0; i < len(c.ManagedOrganizations); i++ {
		for j := i + 1; j < len(c.ManagedOrganizations); j++ {
			orgID1 := c.ManagedOrganizations[i].ID
			orgID2 := c.ManagedOrganizations[j].ID
			key := composeOrgRandMapKey(orgID1, orgID2)
			if _, ok := c.orgRandMap[key]; ok {
				continue
			}
			randList, err := RandIntList(MaxNumTXInEpoch)
			if err != nil {
				return err
			}
			c.orgRandMap[key] = randList
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

func (c *Committee) epochAuditorRandomness() error {
	// generate randomness for auditors
	// distribute randomness to auditors
	var err error
	c.auditorRandList, err = RandIntList(uint(len(c.ManagedAuditors)))
	return err
}
