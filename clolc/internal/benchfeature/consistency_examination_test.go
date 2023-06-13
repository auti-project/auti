package benchfeature

import (
	"math/rand"
	"testing"
	"time"

	"github.com/auti-project/auti/internal/constants"
	"github.com/auti-project/auti/internal/transaction"
	"go.dedis.ch/kyber/v3"
)

const testNumTXs = constants.MaxNumTXInEpoch

func TestConsistencyExaminationCheck(t *testing.T) {
	// setup
	com, auditors, organizations := generateEntities(2)
	_, err := com.InitializeEpoch(auditors, organizations)
	if err != nil {
		t.Fatal(err)
	}
	// test
	txList1 := make([]*transaction.CLOLCLocalPlain, testNumTXs)
	txList2 := make([]*transaction.CLOLCLocalPlain, testNumTXs)
	for i := 0; i < testNumTXs; i++ {
		// generate a random float64
		amount := rand.Float64()
		integerPart := rand.Int()%10000 - 5000
		amount += float64(integerPart)
		currTimestamp := time.Now().UnixNano()
		tx1, tx2 := transaction.NewPairCLOLCLocalPlain(
			string(organizations[0].ID),
			string(organizations[1].ID),
			amount,
			currTimestamp,
		)
		txList1[i] = tx1
		txList2[i] = tx2
	}
	// compute hidden transactions
	var (
		hiddenTXs1   = make([]*transaction.CLOLCLocalHidden, testNumTXs)
		hiddenTXs2   = make([]*transaction.CLOLCLocalHidden, testNumTXs)
		points1      = make([]kyber.Point, testNumTXs)
		points2      = make([]kyber.Point, testNumTXs)
		randScalars1 = make([]kyber.Scalar, testNumTXs)
		randScalars2 = make([]kyber.Scalar, testNumTXs)
	)
	for i := 0; i < testNumTXs; i++ {
		hiddenTX1, point1, scalar1, err := txList1[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs1[i] = hiddenTX1
		points1[i] = point1
		randScalars1[i] = scalar1
		hiddenTX2, point2, scalar2, err := txList2[i].Hide()
		if err != nil {
			t.Fatal(err)
		}
		hiddenTXs2[i] = hiddenTX2
		points2[i] = point2
		randScalars2[i] = scalar2
	}
	res1, err := auditors[0].AccumulateCommitments(organizations[0].ID, hiddenTXs1)
	if err != nil {
		t.Fatal(err)
	}
	res2, err := auditors[1].AccumulateCommitments(organizations[1].ID, hiddenTXs2)
	if err != nil {
		t.Fatal(err)
	}
	txRandList := auditors[0].GetEpochTXRandomness(organizations[0].ID, organizations[1].ID)
	if txRandList == nil {
		t.Fatal("txRandList is nil")
	}
	b1, err := auditors[0].ComputeB(randScalars1, txRandList)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := auditors[1].ComputeB(randScalars2, txRandList)
	if err != nil {
		t.Fatal(err)
	}
	if !auditors[0].CheckResultConsistency(res1, b1, res2, b2) {
		t.Fatal("result consistency check failed")
	}
}
