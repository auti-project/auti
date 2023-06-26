package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/internal/committee/clolc"
	clolc2 "github.com/auti-project/auti/internal/organization/clolc"

	"github.com/auti-project/auti/clolc/internal/timecounter"
	"github.com/auti-project/auti/internal/auditor"
)

func generateEntities(numOrganizations int) (*clolc.Committee, []*auditor.CLOLCAuditor, []*clolc2.Organization) {
	organizations := make([]*clolc2.Organization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = clolc2.New("org" + string(rune(i)))
	}
	auditors := make([]*auditor.CLOLCAuditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = auditor.CLOLCNew("aud"+string(rune(i)), []*clolc2.Organization{organizations[i]})
	}
	com := clolc.New("com", auditors)
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
