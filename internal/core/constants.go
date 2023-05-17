package core

import (
	"math/big"

	"go.dedis.ch/kyber/v3/group/edwards25519"
)

const (
	SecurityParameter uint = 256
	MaxNumTXInEpoch   uint = 1024
)

var (
	Big1             = big.NewInt(1)
	RandIntLimit     = new(big.Int).Lsh(Big1, SecurityParameter)
	kyberSuite       = edwards25519.NewBlakeSHA256Ed25519()
	maxAmountByteLen = kyberSuite.Point().EmbedLen()
	G                = kyberSuite.Point().Base()
	hScalarBytes     = []byte{88, 110, 203, 46, 52, 29, 230, 201, 240, 164, 50, 0, 116, 207, 45, 187, 223, 113, 166, 40,
		12, 27, 15, 50, 235, 140, 55, 192, 37, 22, 130, 239}
	hScalar = kyberSuite.Scalar().SetBytes(hScalarBytes)
	H       = kyberSuite.Point().Base().Mul(hScalar, nil)
)
