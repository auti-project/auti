package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/benchmark/internal/localchain"
	"github.com/auti-project/auti/internal/auditor"
	"github.com/auti-project/auti/internal/committee"
	"github.com/auti-project/auti/internal/organization"
	"go.dedis.ch/kyber/v3"
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

func InitializeEpoch(numOrganizations, iterations int) {
	fmt.Println("CLOLC initialize epoch")
	fmt.Println("Number of organizations:", numOrganizations)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		startTime := time.Now()
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			panic(err)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
}

func TransactionRecordSubmitTX(numTXs, iterations int) error {
	fmt.Println("CLOLC transaction record submit transaction")
	fmt.Println("Number of transactions:", numTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, err := localchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordReadTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record read transaction")
	fmt.Println("Number of transactions:", numTotalTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	txIDs, err := localchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = localchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.BenchReadTX(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record read all transactions")
	fmt.Println("Number of transactions:", numTotalTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	txIDs, err := localchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = localchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := localchain.ReadAllTXs(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordCommitment(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record commitment")
	fmt.Println("Number of transactions:", numTotalTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	for i := 0; i < iterations; i++ {
		dummyTXs, err := localchain.DummyOnChainTransactions(numTotalTXs)
		if err != nil {
			return err
		}
		startTime := time.Now()
		for _, tx := range dummyTXs {
			_, err := tx.ToHidden()
			if err != nil {
				return err
			}
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func TransactionRecordAccumulate(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC transaction record accumulate")
	fmt.Println("Number of transactions:", numTotalTXs)
	fmt.Println("Number of iterations:", iterations)
	fmt.Println("Elapsed times (ms):")
	for i := 0; i < iterations; i++ {
		commitments := make([]kyber.Point, numTotalTXs)
		for i := 0; i < numTotalTXs; i++ {
			plainTX, err := localchain.DummyPlainTransaction()
			if err != nil {
				return err
			}
			_, commitments[i], err = plainTX.Hide()
		}
		kyberSuite := edwards25519.NewBlakeSHA256Ed25519()
		accumulator := kyberSuite.Point().Null()
		startTime := time.Now()
		for _, commitment := range commitments {
			accumulator = accumulator.Add(accumulator, commitment)
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Println(elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}
