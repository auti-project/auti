package benchfeature

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	mt "github.com/txaty/go-merkletree"

	"github.com/auti-project/auti/benchmark/closc/internal/localchain"
	"github.com/auti-project/auti/benchmark/closc/internal/localchaincommit"
	"github.com/auti-project/auti/benchmark/closc/internal/orgchain"
	"github.com/auti-project/auti/benchmark/timecounter"
	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/crypto"
)

const hashByteLen = 32

var numCPU = runtime.NumCPU()

type randInput struct {
	amount       int64
	timestamp    int64
	receiverHash []byte
	counter      uint64
}

func generateRandInputs(num int) []randInput {
	results := make([]randInput, num)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < num; j += step {
				randBytes := make([]byte, constants.SecurityParameterBytes)
				_, err := crand.Read(randBytes)
				if err != nil {
					panic(err)
				}
				results[j] = randInput{
					amount:       rand.Int63(),
					timestamp:    time.Now().UnixNano(),
					receiverHash: randBytes,
					counter:      uint64(j),
				}
			}
		}(i, numCPU)
	}
	wg.Wait()
	return results
}

func TransactionRecordCommitment(num, iterations int) error {
	fmt.Println("[CLOSC-TV] Commitment")
	fmt.Printf("Num: %d, Num iter: %d\n", num, iterations)
	for i := 0; i < iterations; i++ {
		randInputs := generateRandInputs(num)
		startTime := time.Now()
		for j := 0; j < num; j++ {
			if _, _, err := crypto.PedersonCommitWithHash(
				randInputs[j].amount,
				randInputs[j].timestamp,
				randInputs[j].receiverHash,
				randInputs[j].counter,
			); err != nil {
				return err
			}
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

type dummyDataBlock struct {
	data []byte
}

func (d dummyDataBlock) Serialize() ([]byte, error) {
	return d.data, nil
}

func generateDataBlocks(num int) []mt.DataBlock {
	results := make([]mt.DataBlock, num)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < num; j += step {
				randBytes := make([]byte, hashByteLen)
				_, err := crand.Read(randBytes)
				if err != nil {
					panic(err)
				}
				results[j] = dummyDataBlock{data: randBytes}
			}
		}(i, numCPU)
	}
	wg.Wait()
	return results
}

func TransactionRecordMerkleProofGen(depth, iterations int) error {
	fmt.Println("[CLOSC-TV] Merkle proof generation")
	fmt.Printf("Depth: %d, Num iter: %d\n", depth, iterations)
	numDataBlock := 1 << depth
	dataBlocks := generateDataBlocks(numDataBlock)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		_, err := crypto.GenerateMerkleProofs(dataBlocks)
		if err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOSC-TV] Local submit transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := localchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOSC-TV] Local prepare transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTXs)
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
	fmt.Println("[CLOSC-TV] Local read transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOSC-TV] Local read all transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalCommitmentSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOSC-TV] Local submit commitment transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := localchaincommit.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalCommitmentPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOSC-TV] Local prepare commitment transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTXs)
	txIDs, err := localchaincommit.SubmitTX(numTotalTXs)
	if err != nil {
		return err
	}
	if err = localchaincommit.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalCommitmentReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOSC-TV] Local read commitment transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchaincommit.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordLocalCommitmentReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOSC-TV] Local read all commitment transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := localchaincommit.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgSubmitTX(numTXs, iterations int) error {
	fmt.Println("[CLOSC-TV] Org submit transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	for i := 0; i < iterations; i++ {
		_, err := orgchain.SubmitTX(numTXs)
		if err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgPrepareTX(numTotalTXs int) error {
	fmt.Println("[CLOSC-TV] Org prepare transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTXs)
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

func TransactionRecordOrgReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOSC-TV] Org read transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func TransactionRecordOrgReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOSC-TV] Org read all transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := orgchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}
