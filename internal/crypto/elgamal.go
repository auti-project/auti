package crypto

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

var (
	KyberSuite       = edwards25519.NewBlakeSHA256Ed25519()
	maxAmountByteLen = KyberSuite.Point().EmbedLen()
	PointG           = KyberSuite.Point().Base()
	hScalarBytes     = []byte{
		88, 110, 203, 46, 52, 29, 230, 201, 240, 164, 50, 0,
		116, 207, 45, 187, 223, 113, 166, 40, 12, 27, 15, 50,
		235, 140, 55, 192, 37, 22, 130, 239,
	}
	hScalar = KyberSuite.Scalar().SetBytes(hScalarBytes)
	PointH  = KyberSuite.Point().Base().Mul(hScalar, nil)
)

type TypePublicKey kyber.Point
type TypePrivateKey kyber.Scalar

func KeyGen() (privateKey TypePrivateKey, publicKey TypePublicKey, err error) {
	privateKey = KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
	publicKey = KyberSuite.Point().Mul(privateKey, nil)
	return
}

type KeyPair struct {
	PrivateKey TypePrivateKey `json:"private_key"`
	PublicKey  TypePublicKey  `json:"public_key"`
}

type CipherText struct {
	C1 kyber.Point `json:"c_1"`
	C2 kyber.Point `json:"c_2"`
}

func (c *CipherText) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

func DeserializeCipherText(data []byte) (*CipherText, error) {
	c := &CipherText{}
	err := json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
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
	amountPoint := KyberSuite.Point().Embed(amountBytes, KyberSuite.RandomStream())
	randomScalar := KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
	c1 := KyberSuite.Point().Mul(randomScalar, nil)
	c2 := KyberSuite.Point().Add(amountPoint, KyberSuite.Point().Mul(randomScalar, publicKey))
	return &CipherText{c1, c2}, nil
}

func EncryptPoint(publicKey, data kyber.Point) (*CipherText, error) {
	randomScalar := KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
	c1 := KyberSuite.Point().Mul(randomScalar, nil)
	c2 := KyberSuite.Point().Add(data, KyberSuite.Point().Mul(randomScalar, publicKey))
	return &CipherText{c1, c2}, nil
}

func Decrypt(privateKey kyber.Scalar, cipherText *CipherText) (int64, error) {
	amountPoint := KyberSuite.Point().Mul(privateKey, cipherText.C1)
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
