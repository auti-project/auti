package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/benchmark/internal/audchain"
	"github.com/auti-project/auti/clolc/benchmark/internal/localchain"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/organization"
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

func ConsistencyExaminationComputeD(numOrganizations, iterations int) error {
	fmt.Println("CLOLC consistency examination compute D")
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
		_ = auditors[0].ComputeD(randPoint1, randPoint2)
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

func ConsistencyExaminationEncrypt(numOrganizations, iterations int) error {
	fmt.Println("CLOLC consistency examination encrypt")
	fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		counterPartyHashStr := organization.IDHashString(organizations[1].ID)
		_, publicKey, err := crypto.KeyGen()
		if err != nil {
			return err
		}
		randPoint1 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint2 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint3 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		randPoint4 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		startTime := time.Now()
		if _, err := auditors[0].EncryptConsistencyExamResult(
			organizations[0].ID, counterPartyHashStr, randPoint1, randPoint2, randPoint3, randPoint4, publicKey,
		); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		if elapsed.Milliseconds() <= 1 {
			fmt.Printf("Elapsed time: %d ns\n", elapsed.Nanoseconds())
		} else {
			fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
		}
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudSubmitTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC consistency examination submit TX")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, err := audchain.SubmitTX(numTotalTXs)
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

func ConsistencyExaminationAudReadTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC consistency examination read TX")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := audchain.ReadTX(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC consistency examination read all TXs")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		if err := audchain.ReadAllTXs(); err != nil {
			return err
		}
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		fmt.Printf("Elapsed time: %d ms\n", elapsed.Milliseconds())
	}
	fmt.Println()
	return nil
}

func PrepareAudTX(numTotalTXs int) error {
	fmt.Println("CLOLC prepare aud transaction")
	fmt.Printf("Num TX: %d\n", numTotalTXs)
	txIDs, err := audchain.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = audchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}
