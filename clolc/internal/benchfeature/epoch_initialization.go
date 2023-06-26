package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/internal/timecounter"
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/organization"
)

func generateEntities(numOrganizations int) (*committee.CLOLCCommittee, []*auditor.CLOLCAuditor, []*organization.CLOLCOrganization) {
	organizations := make([]*organization.CLOLCOrganization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = organization.CLOLCNew("org" + string(rune(i)))
	}
	auditors := make([]*auditor.CLOLCAuditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = auditor.CLOLCNew("aud"+string(rune(i)), []*organization.CLOLCOrganization{organizations[i]})
	}
	com := committee.CLOLCNew("com", auditors)
	return com, auditors, organizations
}

func InitializeEpoch(numOrganizations, iterations int) error {
	fmt.Println("CLOLC initialize epoch")
	fmt.Printf("Num Org: %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		startTime := time.Now()
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
