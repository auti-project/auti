package audchain

import (
	crand "crypto/rand"
	"runtime"
	"sync"

	"github.com/auti-project/auti/internal/transaction/clolc"

	"go.dedis.ch/kyber/v3/group/edwards25519"

	"github.com/auti-project/auti/internal/crypto"
)

var (
	numCPUs    = runtime.NumCPU()
	kyberSuite = edwards25519.NewBlakeSHA256Ed25519()
)

func DummyOnChainTransactions(numTXs int) []*clolc.AudOnChain {
	results := make([]*clolc.AudOnChain, numTXs)
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

func DummyOnChainTransaction() (*clolc.AudOnChain, error) {
	randIDBytes := make([]byte, 32)
	_, err := crand.Read(randIDBytes)
	if err != nil {
		return nil, err
	}
	randCipherBytes := make([][]byte, 4)
	for i := 0; i < 4; i++ {
		ct := crypto.CipherText{
			C1: kyberSuite.Point().Pick(kyberSuite.RandomStream()),
			C2: kyberSuite.Point().Pick(kyberSuite.RandomStream()),
		}
		randCipherBytes[i], err = ct.Serialize()
		if err != nil {
			return nil, err
		}
	}
	tx := clolc.NewAudPlain(
		randIDBytes, randCipherBytes[0], randCipherBytes[1], randCipherBytes[2], randCipherBytes[3],
	)
	return tx.ToOnChain(), nil
}
