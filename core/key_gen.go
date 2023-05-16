package core

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

var kyberSuite = edwards25519.NewBlakeSHA256Ed25519()

func KeyGen() (privateKey kyber.Scalar, publicKey kyber.Point, err error) {
	privateKey = kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	publicKey = kyberSuite.Point().Mul(privateKey, nil)
	return
}
