package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/auti-project/auti/clolc/application/ledger"
	"github.com/auti-project/auti/internal/transaction"
)

const (
	numTXs      = 90000
	txThreshold = 10000
	trail       = 20
)

func main() {
	submitTX()
	readTX()
}

func readTX() {
	f, err := os.Open("tx_id.log")
	handleErr(err)
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var txIDList []string

	for fileScanner.Scan() {
		txIDList = append(txIDList, fileScanner.Text())
	}

	err = f.Close()
	handleErr(err)
	lc := ledger.NewController()
	defer lc.Close()
	for i := 0; i < trail; i++ {
		idx := rand.Int() % len(txIDList)
		startTime := time.Now()
		_, err := lc.ReadTX(txIDList[idx])
		handleErr(err)
		duration := time.Since(startTime)
		fmt.Println(duration.Milliseconds())
	}
}

func submitTX() {
	lc := ledger.NewController()
	defer lc.Close()
	dummyTXs := dummyTransactions(numTXs)
	for i := 0; i < numTXs; i += txThreshold {
		right := i + txThreshold
		if right > numTXs {
			right = numTXs
		}
		txIDs, err := lc.SubmitBatchTXs(dummyTXs[i:right])
		handleErr(err)
		f, err := os.OpenFile("tx_id.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		handleErr(err)
		for _, id := range txIDs {
			_, err = f.WriteString(id + "\n")
			handleErr(err)
		}
		f.Close()
	}
}

func dummyTransactions(num int) []*transaction.CLOLCOnChain {
	digests := make([]*transaction.CLOLCOnChain, num)
	for i := 0; i < num; i++ {
		digests[i] = dummyTransaction()
	}
	return digests
}

func dummyTransaction() *transaction.CLOLCOnChain {
	currTimestamp := time.Now().UnixNano()
	randAmount := rand.Float64()
	plainTX := transaction.NewCLOLCPlain(
		"Org2",
		randAmount,
		currTimestamp,
	)
	hiddenTX, err := plainTX.Hide()
	handleErr(err)
	return hiddenTX.ToOnChain()
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
