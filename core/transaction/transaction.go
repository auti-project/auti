package transaction

const (
	amountAmplifier = 100
)

type CLOLCPlain struct {
	CounterParty string
	Amount       int64
}

func NewCLOLCPlain(counterParty string, amount float64) *CLOLCPlain {
	amountInt := int64(amount * amountAmplifier)
	return &CLOLCPlain{
		CounterParty: counterParty,
		Amount:       amountInt,
	}
}

type CLOLCCipher struct {
	CounterParty []byte
	Commitment   []byte
}

func NewCLOLCCipher(counterParty, commitment []byte) *CLOLCCipher {
	return &CLOLCCipher{
		CounterParty: counterParty,
		Commitment:   commitment,
	}
}

type CLOLCOnChain struct {
	CounterParty string
	Commitment   string
}

func NewCLOLCOnChain(counterParty, commitment string) *CLOLCOnChain {
	return &CLOLCOnChain{
		CounterParty: counterParty,
		Commitment:   commitment,
	}
}
