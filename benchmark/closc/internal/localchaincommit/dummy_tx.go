package localchaincommit

import (
	crand "crypto/rand"
	"runtime"
	"sync"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/transaction/closc"
)

var numCPUs = runtime.NumCPU()

func DummyCommitmentOnChainTransactions(numTXs int) []*closc.LocalCommitmentOnChain {
	results := make([]*closc.LocalCommitmentOnChain, numTXs)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			for j := idx; j < numTXs; j += step {
				dummyTX, err := DummyCommitmentOnChainTransaction()
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

func DummyCommitmentOnChainTransaction() (*closc.LocalCommitmentOnChain, error) {
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX, err := DummyCommitmentPlainTransaction()
	if err != nil {
		return nil, err
	}
	return plainTX.ToOnChain(), nil
}

func DummyCommitmentPlainTransactions(numTXs int) []*closc.LocalCommitmentPlain {
	results := make([]*closc.LocalCommitmentPlain, numTXs)
	wg := sync.WaitGroup{}
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(idx, step int) {
			defer wg.Done()
			var err error
			for j := idx; j < numTXs; j += step {
				results[j], err = DummyCommitmentPlainTransaction()
				if err != nil {
					panic(err)
				}
			}
		}(i, numCPUs)
	}
	wg.Wait()
	return results
}

func DummyCommitmentPlainTransaction() (*closc.LocalCommitmentPlain, error) {
	randScalar := crypto.KyberSuite.Scalar().Pick(crypto.KyberSuite.RandomStream())
	randPoint := crypto.KyberSuite.Point().Mul(randScalar, nil)
	dummyCommitment, err := randPoint.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return closc.NewLocalCommitmentPlain(dummyCommitment), nil
}
