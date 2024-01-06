package task

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"go.dedis.ch/kyber/v3"

	"github.com/auti-project/auti/benchmark/clolc/internal/blockchain/audchain"
	"github.com/auti-project/auti/benchmark/clolc/internal/blockchain/orgchain"
	"github.com/auti-project/auti/benchmark/timecounter"
	"github.com/auti-project/auti/internal/crypto"
)

func RVVerifyOrgAndAudResult(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-RV] Verify org and aud result")
	for i := 0; i < iterations; i++ {
		fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		dummyOrgChainTX, err := orgchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		dummyAudChainTX, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		startTime := time.Now()
		if _, err = com.VerifyOrgAndAudResult(
			organizations[0].ID,
			auditors[0].ID,
			dummyOrgChainTX,
			dummyAudChainTX,
		); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func RVVerifyAuditPairResult(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-RV] Verify audit pair result")
	for i := 0; i < iterations; i++ {
		fmt.Printf("Num org %d, Num iter: %d\n", numOrganizations, iterations)
		com, auditors, organizations := generateEntities(numOrganizations)
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		dummyTX1, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		dummyTX2, err := audchain.DummyOnChainTransaction()
		if err != nil {
			return err
		}
		startTime := time.Now()
		if _, err = com.VerifyAuditPairResult(
			organizations[0].ID,
			organizations[1].ID,
			auditors[0].ID,
			auditors[1].ID,
			dummyTX1,
			dummyTX2,
		); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func RVBatchDecrypt(iterations, numRoutines int) error {
	if numRoutines <= 0 {
		numRoutines = runtime.NumCPU()
	}
	const numTXs = 255 * 256
	fmt.Println("[CLOLC-RV] Batch decrypt")
	for iter := 0; iter < iterations; iter++ {
		fmt.Printf("Num iter: %d, Num routines: %d\n", iter, numRoutines)
		dummyAudOnChainTXs := audchain.DummyOnChainTransactions(numTXs)
		priKey, _, err := crypto.KeyGen()
		if err != nil {
			return err
		}
		com, _, _ := generateEntities(1)
		runtime.GC()
		startTime := time.Now()
		var wg sync.WaitGroup
		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := idx; j < numTXs; j += numRoutines {
					if _, _, _, err := com.DecryptAudTX(priKey, dummyAudOnChainTXs[j]); err != nil {
						panic(err)
					}
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func RVBatchCheckOrgAndAudPair(iterations, numRoutines int) error {
	if numRoutines <= 0 {
		numRoutines = runtime.NumCPU()
	}
	const numTXs = 256
	fmt.Println("[CLOLC-RV] Batch check org and aud pair")
	for iter := 0; iter < iterations; iter++ {
		fmt.Printf("Num iter: %d, Num routines: %d\n", iter, numRoutines)
		dummyOrgOnChainTXs := orgchain.DummyOnChainTransactions(numTXs)
		dummyPointBList := make([]kyber.Point, numTXs)
		dummyPointCList := make([]kyber.Point, numTXs)
		for i := 0; i < numTXs; i++ {
			dummyPointBList[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			dummyPointCList[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		}
		com, auditors, organizations := generateEntities(1)
		runtime.GC()
		startTime := time.Now()
		var wg sync.WaitGroup
		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := idx; j < numTXs; j += numRoutines {
					if _, err := com.CheckOrgAndAudPair(
						organizations[0].ID,
						auditors[0].ID,
						dummyOrgOnChainTXs[j],
						dummyPointBList[j],
						dummyPointCList[j],
					); err != nil {
						panic(err)
					}
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func RVBatchCheckAudPair(iterations, numRoutines int) error {
	if numRoutines <= 0 {
		numRoutines = runtime.NumCPU()
	}
	const numTXs = 255 * 256
	fmt.Println("[CLOLC-RV] Batch check audit pair")
	for iter := 0; iter < iterations; iter++ {
		fmt.Printf("Num iter: %d, Num routines: %d\n", iter, numRoutines)
		dummyPointC1List := make([]kyber.Point, numTXs)
		dummyPointC2List := make([]kyber.Point, numTXs)
		dummyPointD1List := make([]kyber.Point, numTXs)
		dummyPointD2List := make([]kyber.Point, numTXs)
		for i := 0; i < numTXs; i++ {
			dummyPointC1List[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			dummyPointC2List[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			dummyPointD1List[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			dummyPointD2List[i] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		}
		com, auditors, organizations := generateEntities(2)
		if _, err := com.InitializeEpoch(auditors, organizations); err != nil {
			return err
		}
		runtime.GC()
		startTime := time.Now()
		var wg sync.WaitGroup
		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := idx; j < numTXs; j += numRoutines {
					if _, err := com.CheckAuditPair(
						auditors[0].ID,
						auditors[1].ID,
						dummyPointC1List[j],
						dummyPointC2List[j],
						dummyPointD1List[j],
						dummyPointD2List[j],
					); err != nil {
						panic(err)
					}
				}
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
