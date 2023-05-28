package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/benchmark/internal/localchain"
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/organization"
)

func ConsistencyExaminationAccumulateCommitment(numOrganizations, iterations int) error {
	fmt.Println("CLOLC consistency examination accumulate commitment")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		organizations := make([]*organization.Organization, numOrganizations)
		for i := 0; i < numOrganizations; i++ {
			organizations[i] = organization.New("org" + string(rune(i)))
		}
		auditors := make([]*auditor.Auditor, numOrganizations)
		for i := 0; i < numOrganizations; i++ {
			auditors[i] = auditor.New("aud"+string(rune(i)), []*organization.Organization{organizations[i]})
		}
		com := committee.New("com", auditors)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		dummyTXs := localchain.DummyHiddenTXWithCounterPartyID(organizations[1].ID, int(constants.MaxNumTXInEpoch))
		startTime := time.Now()
		if _, err = auditors[0].AccumulateCommitments(organizations[0].ID, dummyTXs); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		if elapsed.Milliseconds() == 0 {
			fmt.Printf("Elapsed time: %d ns\n", elapsed.Nanoseconds())
		} else {
			fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
		}
	}
	fmt.Println()
	return nil
}
