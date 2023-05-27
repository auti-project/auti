package localchain

import (
	crand "crypto/rand"
	"math/rand"
	"time"

	"github.com/auti-project/auti/internal/transaction"
)

func DummyOnChainTransactions(num int) ([]*transaction.CLOLCOnChain, error) {
	onChainTransactions := make([]*transaction.CLOLCOnChain, num)
	var err error
	for i := 0; i < num; i++ {
		onChainTransactions[i], err = DummyOnChainTransaction()
		if err != nil {
			return nil, err
		}
	}
	return onChainTransactions, nil
}

func DummyOnChainTransaction() (*transaction.CLOLCOnChain, error) {
	dummyCounterPartyBytes := make([]byte, 32)
	_, err := crand.Read(dummyCounterPartyBytes)
	if err != nil {
		return nil, err
	}
	plainTX, err := DummyPlainTransaction()
	if err != nil {
		return nil, err
	}
	hiddenTX, _, err := plainTX.Hide()
	if err != nil {
		return nil, err
	}
	return hiddenTX.ToOnChain(), nil
}

func DummyPlainTransaction() (*transaction.CLOLCPlain, error) {
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
	return plainTX, nil
}
