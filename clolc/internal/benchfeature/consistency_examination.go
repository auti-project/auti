package benchfeature

import (
	"fmt"
	"time"

	"github.com/auti-project/auti/clolc/internal/audchain"
	"github.com/auti-project/auti/clolc/internal/localchain"
	"github.com/auti-project/auti/clolc/internal/timecounter"
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
		dummyTXs := localchain.DummyHiddenTXWithCounterPartyID(organizations[1].ID, constants.MaxNumTXInEpoch)
		startTime := time.Now()
		if _, err = auditors[0].AccumulateCommitments(organizations[0].ID, dummyTXs); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
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
		for i := 0; i < constants.MaxNumTXInEpoch; i++ {
			randScalars1[i] = crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
			randScalars2[i] = crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
		}
		startTime := time.Now()
		if _, err = auditors[0].ComputeB(randScalars1, randScalars2); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
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
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
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
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
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
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudSubmitTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC consistency examination submit TX")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if _, err := audchain.SubmitTX(numTotalTXs); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudReadTX(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC consistency examination read TX")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := audchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("CLOLC consistency examination read all TXs")
	fmt.Printf("Num total TXs %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := audchain.ReadAllTXsByPage(); err != nil {
			return err
		}
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

func ConsistencyExaminationDecrypt(iterations int) error {
	fmt.Println("CLOLC consistency examination decrypt")
	fmt.Printf("Num iter: %d\n", iterations)
	com, auditors, organizations := generateEntities(2)
	_, err := com.InitializeEpoch(auditors, organizations)
	if err != nil {
		return err
	}
	for i := 0; i < iterations; i++ {
		dummyTX, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		orgIDHashStr := organization.IDHashString(organizations[0].ID)
		startTime := time.Now()
		if _, _, err := auditors[0].DecryptResAndB(orgIDHashStr, dummyTX); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationCheck(iterations int) error {
	fmt.Println("CLOLC consistency examination check")
	fmt.Printf("Num iter: %d\n", iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(2)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		randPoints := make([]kyber.Point, 4)
		for i := 0; i < 4; i++ {
			randPoints[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		}
		startTime := time.Now()
		_ = auditors[0].CheckResultConsistency(
			randPoints[0], randPoints[1], randPoints[2], randPoints[3],
		)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
