package transaction

const (
	amountAmplifier = 100
)

type CLOLCPlain struct {
	Sender   string
	Receiver string
	Amount   int64
}

func NewCLOLCPlain(sender, receiver string, amount float64) *CLOLCPlain {
	amountInt := int64(amount * amountAmplifier)
	return &CLOLCPlain{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amountInt,
	}
}

type CLOLCCipher struct {
	Sender     []byte
	Receiver   []byte
	Commitment []byte
}

type CLOLCOnChain struct {
	Sender     string
	Receiver   string
	Commitment string
}
