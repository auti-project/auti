package benchfeature

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	mt "github.com/txaty/go-merkletree"

	"github.com/auti-project/auti/benchmark/timecounter"
	"github.com/auti-project/auti/internal/closc/auditor"
	"github.com/auti-project/auti/internal/closc/transaction"
	"github.com/auti-project/auti/internal/crypto"
)

const mergeTreeDepth = 20

var merkleTreeConfig = &mt.Config{
	DisableLeafHashing: true,
	RunInParallel:      true,
}

func genDummyDataBlockAndProof(treeDepth int) (dataBlocks []mt.DataBlock, merkleProofs []*mt.Proof, err error) {
	numTXs := 1 << treeDepth
	dummyDataBlocks := generateDataBlocks(numTXs)
	tree, err := mt.New(merkleTreeConfig, dummyDataBlocks)
	if err != nil {
		return nil, nil, err
	}
	return dummyDataBlocks, tree.Proofs, nil
}

func genDummyLocalOnChainTX(treeDepth int) (txList []transaction.LocalOnChain, err error) {
	numTXs := 1 << treeDepth
	dummyCommitments, merkleProofs, err := genDummyDataBlockAndProof(treeDepth)
	if err != nil {
		return nil, err
	}
	txList = make([]transaction.LocalOnChain, numTXs)
	for i := 0; i < numTXs; i++ {
		dummyCommitmentByte, err := dummyCommitments[i].Serialize()
		if err != nil {
			return nil, err
		}
		dummyCommitmentStr := hex.EncodeToString(dummyCommitmentByte)
		merkleProofBytes, err := crypto.MerkleProofMarshal(merkleProofs[i])
		if err != nil {
			return nil, err
		}
		merkleProofStr := hex.EncodeToString(merkleProofBytes)
		txList[i] = transaction.LocalOnChain{
			Commitment:  dummyCommitmentStr,
			MerkleProof: merkleProofStr,
		}
	}
	return txList, nil
}

func ConsistencyExaminationMerkleProofVerify(treeDepth, iterations int) error {
	fmt.Println("[CLOLC-CE] Merkle Proof Verify")
	fmt.Printf("Tree depth: %d, Num iter: %d\n", treeDepth, iterations)
	txList, err := genDummyLocalOnChainTX(treeDepth)
	if err != nil {
		return err
	}
	numTXs := 1 << treeDepth
	aud := auditor.New("aud", nil)
	for i := 0; i < iterations; i++ {
		randIdx := rand.Int() % numTXs
		startTime := time.Now()
		ret, err := aud.VerifyMerkleProof(txList[randIdx])
		if err != nil {
			return err
		}
		if ret != 1 {
			return fmt.Errorf("merkle proof verification failed")
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

// randIndexes generates the random indexes without duplication
func randIndexes(numIdx, max int) []int {
	if numIdx > max {
		numIdx = max
	}

	// Generate a pool of indexes
	pool := make([]int, max)
	for i := 0; i < max; i++ {
		pool[i] = i
	}

	// Shuffle the pool using Fisher-Yates algorithm
	for i := max - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		pool[i], pool[j] = pool[j], pool[i]
	}

	// Take the first numIdx elements from the shuffled pool
	return pool[:numIdx]

}

func ConsistencyExaminationMerkleProofMerge(numTXs, iterations int) error {
	fmt.Println("[CLOLC-CE] Merkle Proof Merge")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTXs, iterations)
	numTotalTXs := 1 << mergeTreeDepth
	dummyBlocks, dummyProofs, err := genDummyDataBlockAndProof(mergeTreeDepth)
	aud := auditor.New("aud", nil)
	for i := 0; i < iterations; i++ {
		if err != nil {
			return err
		}
		indexes := randIndexes(numTXs, numTotalTXs)
		selectedBlocks := make([]mt.DataBlock, numTXs)
		selectedProofs := make([]*mt.Proof, numTXs)
		for j := 0; j < numTXs; j++ {
			selectedBlocks[j] = dummyBlocks[indexes[j]]
			selectedProofs[j] = dummyProofs[indexes[j]]
		}
		startTime := time.Now()
		if _, err = aud.MergeProof(selectedBlocks, dummyProofs); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
