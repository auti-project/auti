package localchain

import (
	crand "crypto/rand"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

var (
	numCPUs = runtime.NumCPU()
)

func DummyOnChainTransactions(numTXs int) []*transaction.CLOLCOnChain {
	results := make([]*transaction.CLOLCOnChain, numTXs)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				dummyTX, err := DummyOnChainTransaction()
				if err != nil {
					panic(err)
				}
				results[j] = dummyTX
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyOnChainTransaction() (*transaction.CLOLCOnChain, error) {
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX, err := DummyPlainTransaction()
	if err != nil {
		return nil, err
	}
	hiddenTX, _, err := plainTX.Hide()
	if err != nil {
		return nil, err
	}
	return hiddenTX.ToOnChain(), nil
}

func DummyPlainTransactions(numTXs int) []*transaction.CLOLCPlain {
	results := make([]*transaction.CLOLCPlain, numTXs)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			var err error
			for j := idx; j < numTXs; j += step {
				results[j], err = DummyPlainTransaction()
				if err != nil {
					panic(err)
				}
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyPlainTransaction() (*transaction.CLOLCPlain, error) {
	currTimestamp := time.Now().UnixNano()
	randAmount := rand.Float64()
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX := transaction.NewCLOLCPlain(
		string(dummyCounterPartyBytes),
		randAmount,
		currTimestamp,
	)
	return plainTX, nil
}

func DummyHiddenTXCommitments(numTXs int) []kyber.Point {
	results := make([]kyber.Point, numTXs)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				tx, err := DummyPlainTransaction()
				if err != nil {
					panic(err)
				}
				_, commitment, err := tx.Hide()
				if err != nil {
					panic(err)
				}
				results[j] = commitment
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}
