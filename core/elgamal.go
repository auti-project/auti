package core

import (
	"errors"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

var maxAmountByteLen = kyberSuite.Point().EmbedLen()

type CipherText struct {
	C1 kyber.Point
	C2 kyber.Point
}

func Encrypt(publicKey kyber.Point, amount *big.Int) (*CipherText, error) {
	// Embed the amount into a curve point
	amountPoint := kyberSuite.Point().Embed(amount.Bytes(), kyberSuite.RandomStream())
	if maxAmountByteLen < len(amount.Bytes()) {
		return nil, errors.New("amount is too large")
	}
	randomScalar := kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	c1 := kyberSuite.Point().Mul(randomScalar, nil)
	c2 := kyberSuite.Point().Add(amountPoint, kyberSuite.Point().Mul(randomScalar, publicKey))
	return &CipherText{c1, c2}, nil
}

func Decrypt(privateKey kyber.Scalar, cipherText *CipherText) (*big.Int, error) {
	amountPoint := kyberSuite.Point().Mul(privateKey, cipherText.C1)
	amountPoint.Neg(amountPoint)
	amountPoint.Add(amountPoint, cipherText.C2)
	amountBytes, err := amountPoint.Data()
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(amountBytes), nil
}
