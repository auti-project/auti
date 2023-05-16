package core

import "math/big"

const (
	SecurityParameter uint = 256
	MaxNumTXInEpoch   uint = 1024
	UniformByteLength      = 64
)

var (
	Big1         = big.NewInt(1)
	RandIntLimit = new(big.Int).Lsh(Big1, SecurityParameter)
)
