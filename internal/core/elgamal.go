package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"go.dedis.ch/kyber/v3"
)

func KeyGen() (privateKey kyber.Scalar, publicKey kyber.Point, err error) {
	privateKey = kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	publicKey = kyberSuite.Point().Mul(privateKey, nil)
	return
}

type KeyPair struct {
	PrivateKey kyber.Scalar
	PublicKey  kyber.Point
}

type CipherText struct {
	C1 kyber.Point
	C2 kyber.Point
}

func Encrypt(publicKey kyber.Point, amount int64) (*CipherText, error) {
	// Embed the amount into a curve point
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, err
	}
	if maxAmountByteLen < len(amountBytes) {
		return nil, errors.New("amount is too large")
	}
	amountPoint := kyberSuite.Point().Embed(amountBytes, kyberSuite.RandomStream())
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

func int64ToBytes(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, fmt.Errorf("int64ToBytes: %v", err)
	}
	return buf.Bytes(), nil
}
