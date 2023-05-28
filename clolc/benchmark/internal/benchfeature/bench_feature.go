package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/benchmark/internal/localchain"
	"github.com/auti-project/auti/clolc/benchmark/internal/orgchain"
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/organization"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

func generateEntities(numOrganizations int) (*committee.Committee, []*auditor.Auditor, []*organization.Organization) {
	organizations := make([]*organization.Organization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = organization.New("org" + string(rune(i)))
	}
	auditors := make([]*auditor.Auditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = auditor.New("aud"+string(rune(i)), []*organization.Organization{organizations[i]})
	}
	com := committee.New("com", auditors)
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
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalSubmitTX(numTXs, iterations int) error {
	fmt.Println("CLOLC transaction record local submit transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, err := localchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func PrepareLocalTX(numTotalTXs int) error {
	fmt.Println("CLOLC prepare transaction")
	fmt.Printf("Num TX: %d\n", numTotalTXs)
	txIDs, err := localchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = localchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalReadTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record local read transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.ReadTX(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record local read all transactions")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.ReadAllTXs(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordCommitment(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record commitment")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		dummyTXs := localchain.DummyPlainTransactions(numTotalTXs)
		startTime := time.Now()
		for _, tx := range dummyTXs {
			_, _, err := tx.Hide()
			if err != nil {
				return err
			}
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordAccumulate(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record accumulate")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		dummyCommitments := localchain.DummyHiddenTXCommitments(numTotalTXs)
		kyberSuite := edwards25519.NewBlakeSHA256Ed25519()
		accumulator := kyberSuite.Point().Null()
		startTime := time.Now()
		for _, commitment := range dummyCommitments {
			accumulator = accumulator.Add(accumulator, commitment)
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

func PrepareOrgTX(numTotalTXs int) error {
	fmt.Println("CLOLC prepare transaction")
	fmt.Printf("Num TX: %d\n", numTotalTXs)
	txIDs, err := orgchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = orgchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgSubmitTX(numTXs, iterations int) error {
	fmt.Println("CLOLC transaction record submit transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, err := orgchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgReadTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record read transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := orgchain.ReadTX(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record read all transactions")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := orgchain.ReadAllTXs(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

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
