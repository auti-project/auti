package benchfeature

import (
	"fmt"
	"time"

	"go.dedis.ch/kyber/v3/group/edwards25519"

	"github.com/auti-project/auti/benchmark/clolc/internal/localchain"
	"github.com/auti-project/auti/benchmark/clolc/internal/orgchain"
	"github.com/auti-project/auti/benchmark/timecounter"
)

func TransactionRecordLocalSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Local submit transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := localchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func PrepareLocalTX(numTotalTXs int) error {
	fmt.Println("[CLOLC-TR] Prepare local transaction")
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
	fmt.Println("[CLOLC-TR] Local read transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Local read all transactions")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordCommitment(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Commitment")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		dummyTXs := localchain.DummyPlainTransactions(numTotalTXs)
		startTime := time.Now()
		for _, tx := range dummyTXs {
			if _, _, _, err := tx.Hide(); err != nil {
				return err
			}
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func TransactionRecordAccumulate(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Accumulate")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		dummyCommitments := localchain.DummyHiddenTXCommitments(numTotalTXs)
		kyberSuite := edwards25519.NewBlakeSHA256Ed25519()
		accumulator := kyberSuite.Point().Null()
		startTime := time.Now()
		for _, commitment := range dummyCommitments {
			accumulator = accumulator.Add(accumulator, commitment)
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func PrepareOrgTX(numTotalTXs int) error {
	fmt.Println("[CLOLC-TR] Prepare transaction")
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
	fmt.Println("[CLOLC-TR] Submit transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		if _, err := orgchain.SubmitTX(numTXs); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Read transaction")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgReadAllTXs(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-TR] Read all transactions")
	fmt.Printf("Num TX: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}