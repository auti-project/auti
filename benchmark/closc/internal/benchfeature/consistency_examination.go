package benchfeature

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	mt "github.com/txaty/go-merkletree"
	"go.dedis.ch/kyber/v3"

	"github.com/auti-project/auti/benchmark/closc/internal/audchain"
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

func genDummyDataBlockAndProof(treeDepth int) (
	dataBlocks []mt.DataBlock, merkleProofs []*mt.Proof, merkleRoot []byte, err error,
) {
	numTXs := 1 << treeDepth
	dummyDataBlocks := generateDataBlocks(numTXs)
	tree, err := mt.New(merkleTreeConfig, dummyDataBlocks)
	if err != nil {
		return nil, nil, nil, err
	}
	return dummyDataBlocks, tree.Proofs, tree.Root, nil
}

func genDummyLocalOnChainTX(treeDepth int) (txList []transaction.LocalOnChain, err error) {
	numTXs := 1 << treeDepth
	dummyCommitments, merkleProofs, root, err := genDummyDataBlockAndProof(treeDepth)
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
			MerkleRoot:  hex.EncodeToString(root),
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
		runtime.GC()
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
	dummyBlocks, dummyProofs, _, err := genDummyDataBlockAndProof(mergeTreeDepth)
	if err != nil {
		return err
	}
	aud := auditor.New("aud", nil)
	for i := 0; i < iterations; i++ {
		indexes := randIndexes(numTXs, numTotalTXs)
		selectedBlocks := make([]mt.DataBlock, numTXs)
		selectedProofs := make([]*mt.Proof, numTXs)
		for j := 0; j < numTXs; j++ {
			selectedBlocks[j] = dummyBlocks[indexes[j]]
			selectedProofs[j] = dummyProofs[indexes[j]]
		}
		startTime := time.Now()
		if _, err = aud.MergeProof(selectedBlocks, selectedProofs); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
		runtime.GC()
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationSummarizeMerkleProofVerificationResults(numResults, iterations int) error {
	fmt.Println("[CLOLC-CE] Summarize Merkle Proof Verification Results")
	fmt.Printf("Num results: %d, Num iter: %d\n", numResults, iterations)
	results := make([]uint, numResults)
	for i := 0; i < iterations; i++ {
		for j := 0; j < numResults; j++ {
			results[j] = uint(rand.Int() % 2)
		}
		aud := auditor.New("aud", nil)
		startTime := time.Now()
		aud.SummarizeMerkleProofVerificationResults(results)
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationVerifyCommitments(numCommitments, iterations int) error {
	fmt.Println("[CLOLC-CE] Verify Commitments")
	fmt.Printf("Num commitments: %d, Num iter: %d\n", numCommitments, iterations)
	aud := auditor.New("aud", nil)
	for i := 0; i < iterations; i++ {
		commitments1 := make([][]byte, numCommitments)
		commitments2 := make([][]byte, numCommitments)
		hashPoints1 := make([]kyber.Point, numCommitments)
		hashPoints2 := make([]kyber.Point, numCommitments)
		var err error
		for j := 0; j < numCommitments; j++ {
			randPoint1 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			randPoint2 := crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			commitments1[j], err = randPoint1.MarshalBinary()
			if err != nil {
				return err
			}
			commitments2[j], err = randPoint2.MarshalBinary()
			if err != nil {
				return err
			}
			hashPoints1[j] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
			hashPoints2[j] = crypto.KyberSuite.Point().Pick(crypto.KyberSuite.RandomStream())
		}
		startTime := time.Now()
		if _, err = aud.VerifyCommitments(commitments1, commitments2, hashPoints1, hashPoints2); err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudSubmitTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-CE] Aud Submit transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if _, err := audchain.SubmitTX(numTotalTXs); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudPrepareTX(numTotalTX int) error {
	fmt.Println("[CLOLC-CE] Aud Prepare transaction")
	fmt.Printf("Num TXs: %d\n", numTotalTX)
	txIDs, err := audchain.SubmitTX(numTotalTX)
	if err != nil {
		return err
	}
	if err = audchain.SaveTXIDs(txIDs); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudReadTX(numTotalTXs, iterations int) error {
	fmt.Println("[CLOLC-CE] Aud Read transaction")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotalTXs, iterations)
	for i := 0; i < iterations; i++ {
		if err := audchain.ReadTX(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func ConsistencyExaminationAudReadAllTXs(numTotals, iterations int) error {
	fmt.Println("[CLOLC-CE] Aud Read all transactions")
	fmt.Printf("Num TXs: %d, Num iter: %d\n", numTotals, iterations)
	for i := 0; i < iterations; i++ {
		if err := audchain.ReadAllTXsByPage(); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}
