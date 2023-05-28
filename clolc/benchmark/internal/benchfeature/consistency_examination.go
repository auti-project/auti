package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/benchmark/internal/localchain"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/crypto"
	"go.dedis.ch/kyber/v3"
)

func ConsistencyExaminationAccumulateCommitment(numOrganizations, iterations int) error {
	fmt.Println("CLOLC consistency examination accumulate commitment")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
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

func ConsistencyExaminationComputeB(numOrganizations, iterations int) error {
	fmt.Println("CLOLC consistency examination compute B")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		randScalars1 := make([]kyber.Scalar, constants.MaxNumTXInEpoch)
		randScalars2 := make([]kyber.Scalar, constants.MaxNumTXInEpoch)
		for i := uint(0); i < constants.MaxNumTXInEpoch; i++ {
			randScalars1[i] = crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
			randScalars2[i] = crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
		}
		startTime := time.Now()
		if _, err = auditors[0].ComputeB(randScalars1, randScalars2); err != nil {
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

func ConsistencyExaminationComputeC(numOrganizations, iterations int) error {
	fmt.Println("CLOLC consistency examination compute C")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}

		randPoint1 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint2 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		startTime := time.Now()
		_ = auditors[0].ComputeC(randPoint1, randPoint2)
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
