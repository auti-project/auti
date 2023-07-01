package audchain

import (
	crand "crypto/rand"
	"errors"
	"math/rand"
	"runtime"
	"sync"

	"go.dedis.ch/kyber/v3/group/edwards25519"

	"github.com/auti-project/auti/internal/closc/transaction"
	"github.com/auti-project/auti/internal/crypto"
)

const (
	treeDepth              = 20
	hashByteLen            = 32
	numRandHashBytes       = 1 << 13
	numHashBytesLowerBound = 1 << 7
)

var (
	numCPUs              = runtime.NumCPU()
	kyberSuite           = edwards25519.NewBlakeSHA256Ed25519()
	onceGenRandHashBytes sync.Once
	allRandHashByteList  [][]byte
)

func genAllRandHashBytes() [][]byte {
	onceGenRandHashBytes.Do(func() {
		allRandHashByteList = make([][]byte, numRandHashBytes)
		for i := 0; i < numRandHashBytes; i++ {
			allRandHashByteList[i] = make([]byte, hashByteLen)
			_, err := crand.Read(allRandHashByteList[i])
			if err != nil {
				panic(err)
			}
		}
	})
	return allRandHashByteList
}

func genRandHashBytes(numHashes int) ([][]byte, error) {
	if numHashes < numHashBytesLowerBound {
		return nil, errors.New("numHashes is too small")
	}
	if numHashes > numRandHashBytes {
		return nil, errors.New("numHashes is too large")
	}
	allRandHashBytes := genAllRandHashBytes()
	return allRandHashBytes[:numHashes], nil
}

func genRandBatchProof() (*crypto.MerkleBatchProof, error) {
	randNum := rand.Int() % (numRandHashBytes - numHashBytesLowerBound)
	randNum += numHashBytesLowerBound
	randHashBytes, err := genRandHashBytes(randNum)
	if err != nil {
		return nil, err
	}
	bp := &crypto.MerkleBatchProof{
		Nodes:   make([]crypto.ProofNode, randNum),
		Indexes: make([]int, randNum),
	}
	for i := 0; i < randNum; i++ {
		bp.Nodes[i] = crypto.ProofNode{
			Data: randHashBytes[i],
			Coordinate: [2]int{
				rand.Int() % treeDepth,
				rand.Int() % (1 << uint(treeDepth)),
			},
		}
		bp.Indexes[i] = i
	}
	return bp, nil
}

func DummyOnChainTransactions(numTXs int) []*transaction.AudOnChain {
	results := make([]*transaction.AudOnChain, numTXs)
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

func DummyOnChainTransaction() (*transaction.AudOnChain, error) {
	randIDBytes := make([]byte, 32)
	_, err := crand.Read(randIDBytes)
	if err != nil {
		return nil, err
	}
	randPoint := kyberSuite.Point().Pick(kyberSuite.RandomStream())
	randBatchProof, err := genRandBatchProof()
	if err != nil {
		return nil, err
	}
	tx, err := transaction.NewAudPlainFromPointAndProof(randPoint, randBatchProof)
	if err != nil {
		return nil, err
	}
	return tx.ToOnChain(), nil
}
