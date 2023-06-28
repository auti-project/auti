package orgchain

import (
	"crypto/sha256"
	"runtime"
	"sync"

	"github.com/auti-project/auti/internal/crypto"
	"github.com/auti-project/auti/internal/transaction/clolc"
)

var (
	numCPUs = runtime.NumCPU()
)

func DummyOnChainTransactions(numTXs int) []*clolc.OrgOnChain {
	results := make([]*clolc.OrgOnChain, numTXs)
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

func DummyOnChainTransaction() (*clolc.OrgOnChain, error) {
	randID, err := crypto.RandBytes()
	if err != nil {
		return nil, err
	}
	sha256Func := sha256.New()
	sha256Func.Write(randID)
	randIDHashBytes := sha256Func.Sum(nil)
	suite := crypto.KyberSuite
	randIDScalar := suite.Scalar().SetBytes(randIDHashBytes)
	randScalar := suite.Scalar().Pick(suite.RandomStream())
	randPoint := suite.Point().Mul(randScalar, suite.Point().Base())
	accumulator := suite.Point().Mul(randIDScalar, randPoint)
	accumulatorBytes, err := accumulator.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tx := clolc.NewOrgPlain(accumulatorBytes)
	return tx.ToOnChain(), nil
}
