package localchain

import (
	"bufio"
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/auti-project/auti/internal/transaction"
)

const (
	txThreshold = 10000
	txIDLogPath = "lc_tx_id.log"
)

func BenchReadTX() error {
	f, err := os.Open(txIDLogPath)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var txIDList []string
	for fileScanner.Scan() {
		txIDList = append(txIDList, fileScanner.Text())
	}
	err = f.Close()
	if err != nil {
		return err
	}
	lc, err := NewController(audWalletPath, audWalletLabel, aud1CCPPath)
	if err != nil {
		return err
	}
	defer lc.Close()
	idx := rand.Int() % len(txIDList)
	startTime := time.Now()
	if _, err = lc.ReadTX(txIDList[idx]); err != nil {
		return err
	}
	duration := time.Since(startTime)
	fmt.Println(duration.Milliseconds())

	return nil
}

func BenchSubmitTX(numTXs int) ([]string, error) {
	lc, err := NewController(orgWalletPath, orgWalletLabel, org1CCPPath)
	if err != nil {
		return nil, err
	}
	defer lc.Close()
	dummyTXs, err := dummyTransactions(numTXs)
	if err != nil {
		return nil, err
	}
	var txIDs []string
	for j := 0; j < numTXs; j += txThreshold {
		right := j + txThreshold
		if right > numTXs {
			right = numTXs
		}
		batchTXIDs, err := lc.SubmitBatchTXs(dummyTXs[j:right])
		if err != nil {
			return nil, err
		}
		txIDs = append(txIDs, batchTXIDs...)
	}
	return txIDs, nil
}

func dummyTransactions(num int) ([]*transaction.CLOLCOnChain, error) {
	digests := make([]*transaction.CLOLCOnChain, num)
	var err error
	for i := 0; i < num; i++ {
		digests[i], err = dummyTransaction()
		if err != nil {
			return nil, err
		}
	}
	return digests, nil
}

func dummyTransaction() (*transaction.CLOLCOnChain, error) {
	currTimestamp := time.Now().UnixNano()
	randAmount := rand.Float64()
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX := transaction.NewCLOLCPlain(
		string(dummyCounterPartyBytes),
		randAmount,
		currTimestamp,
	)
	hiddenTX, err := plainTX.Hide()
	if err != nil {
		return nil, err
	}
	return hiddenTX.ToOnChain(), nil
}

func SaveTXIDs(txIDs []string) error {
	f, err := os.OpenFile(txIDLogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	for _, id := range txIDs {
		if _, err = f.WriteString(id + "\n"); err != nil {
			return err
		}
	}
	return nil
}
