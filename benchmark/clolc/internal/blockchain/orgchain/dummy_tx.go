package orgchain

import (
	"crypto/sha256"
	"runtime"
	"sync"

	"github.com/auti-project/auti/internal/clolc/transaction"
	"github.com/auti-project/auti/internal/crypto"
)

var (
	numCPUs = runtime.NumCPU()
)

func DummyPlainTransactions(numTXs int) []*transaction.OrgPlain {
	results := make([]*transaction.OrgPlain, numTXs)
	var wg sync.WaitGroup
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

func DummyPlainTransaction() (*transaction.OrgPlain, error) {
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
	tx := transaction.NewOrgPlain(accumulatorBytes)
	return tx, nil
}

func DummyOnChainTransactions(numTXs int) []*transaction.OrgOnChain {
	results := make([]*transaction.OrgOnChain, numTXs)
	var wg sync.WaitGroup
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

func DummyOnChainTransaction() (*transaction.OrgOnChain, error) {
	plainTX, err := DummyPlainTransaction()
	if err != nil {
		return nil, err
	}
	return plainTX.ToOnChain(), nil
}
