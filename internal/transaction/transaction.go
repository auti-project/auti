package transaction

const (
	amountAmplifier = 100
)

type CLOLCPlain struct {
	CounterParty string
	Amount       int64
	Timestamp    int64
}

func NewCLOLCPlain(counterParty string, amount, timestamp float64) *CLOLCPlain {
	amountInt := int64(amount * amountAmplifier)
	return &CLOLCPlain{
		CounterParty: counterParty,
		Amount:       amountInt,
	}
}

type CLOLCCipher struct {
	CounterParty []byte
	Commitment   []byte
	Timestamp    int64
}

func NewCLOLCCipher(counterParty, commitment []byte, timestamp int64) *CLOLCCipher {
	return &CLOLCCipher{
		CounterParty: counterParty,
		Commitment:   commitment,
	}
}

type CLOLCOnChain struct {
	CounterParty string
	Commitment   string
	Timestamp    string
}

func NewCLOLCOnChain(counterParty, commitment, timestamp string) *CLOLCOnChain {
	return &CLOLCOnChain{
		CounterParty: counterParty,
		Commitment:   commitment,
	}
}
