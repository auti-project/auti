package transaction

import "encoding/hex"

type CLOLCAudPlain struct {
	Accumulator []byte
}

func NewCLOLCAudPlain(accumulator []byte) *CLOLCAudPlain {
	return &CLOLCAudPlain{
		Accumulator: accumulator,
	}
}

func (c *CLOLCAudPlain) ToOnChain() *CLOLCAudOnChain {
	accumulatorString := hex.EncodeToString(c.Accumulator)
	return NewCLOLCAudOnChain(accumulatorString)
}

type CLOLCAudOnChain struct {
	Accumulator string `json:"accumulator"`
}

func NewCLOLCAudOnChain(accumulator string) *CLOLCAudOnChain {
	return &CLOLCAudOnChain{
		Accumulator: accumulator,
	}
}

func (c *CLOLCAudOnChain) ToPlain() (*CLOLCAudPlain, error) {
	accumulatorBytes, err := hex.DecodeString(c.Accumulator)
	if err != nil {
		return nil, err
	}
	return NewCLOLCAudPlain(accumulatorBytes), nil
}

//func (c *CLOLCAudOnChain) KeyVal()
