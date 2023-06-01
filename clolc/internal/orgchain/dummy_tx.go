package orgchain

import (
	crand "crypto/rand"
	"crypto/sha256"
	"math/big"
	"runtime"
	"sync"

	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

var (
	numCPUs     = runtime.NumCPU()
	big1        = big.NewInt(1)
	randIDLimit = new(big.Int).Lsh(big1, constants.SecurityParameter)
	kyberSuite  = edwards25519.NewBlakeSHA256Ed25519()
)

func DummyOnChainTransactions(numTXs int) []*transaction.CLOLCOrgOnChain {
	results := make([]*transaction.CLOLCOrgOnChain, numTXs)
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

func DummyOnChainTransaction() (*transaction.CLOLCOrgOnChain, error) {
	randID, err := crand.Int(crand.Reader, randIDLimit)
	if err != nil {
		return nil, err
	}
	randIDBytes := randID.Bytes()
	sha256Func := sha256.New()
	sha256Func.Write(randIDBytes)
	randIDHashBytes := sha256Func.Sum(nil)
	randIDScalar := kyberSuite.Scalar().SetBytes(randIDHashBytes)
	randScalar := kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	randPoint := kyberSuite.Point().Mul(randScalar, kyberSuite.Point().Base())
	accumulator := kyberSuite.Point().Mul(randIDScalar, randPoint)
	accumulatorBytes, err := accumulator.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tx := transaction.NewCLOLCOrgPlain(accumulatorBytes)
	return tx.ToOnChain(), nil
}
