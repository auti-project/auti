package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"go.dedis.ch/kyber/v3"
)

type TypePublicKey kyber.Point
type TypeSecretKey kyber.Scalar

func KeyGen() (privateKey TypeSecretKey, publicKey TypePublicKey, err error) {
	privateKey = kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	publicKey = kyberSuite.Point().Mul(privateKey, nil)
	return
}

type KeyPair struct {
	PrivateKey TypeSecretKey
	PublicKey  TypePublicKey
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

func Decrypt(privateKey kyber.Scalar, cipherText *CipherText) (int64, error) {
	amountPoint := kyberSuite.Point().Mul(privateKey, cipherText.C1)
	amountPoint.Neg(amountPoint)
	amountPoint.Add(amountPoint, cipherText.C2)
	amountBytes, err := amountPoint.Data()
	if err != nil {
		return 0, err
	}
	return bytesToInt64(amountBytes)
}

func int64ToBytes(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, fmt.Errorf("int64ToBytes: %v", err)
	}
	return buf.Bytes(), nil
}

func bytesToInt64(b []byte) (int64, error) {
	buf := bytes.NewReader(b)
	var i int64
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, fmt.Errorf("bytesToInt64: %v", err)
	}
	return i, nil
}
