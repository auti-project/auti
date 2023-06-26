package auditor

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/auti-project/auti/internal/crypto"
	"go.dedis.ch/kyber/v3"
)

var (
	sha256Func = sha256.New()
)

type TypeID string
type TypeEpochID []byte

func IDHashBytes(id TypeID) []byte {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	return sha256Func.Sum(nil)
}

func IDHashString(id TypeID) string {
	return hex.EncodeToString(IDHashBytes(id))
}

func IDHashScalar(id TypeID) kyber.Scalar {
	return crypto.KyberSuite.Scalar().SetBytes(IDHashBytes(id))
}

func IDHashPoint(id TypeID) kyber.Point {
	return crypto.KyberSuite.Point().Mul(IDHashScalar(id), nil)
}

func EpochIDHashBytes(epochID TypeEpochID) []byte {
	defer sha256Func.Reset()
	sha256Func.Write(epochID)
	return sha256Func.Sum(nil)
}

func EpochIDHashString(epochID TypeEpochID) string {
	return hex.EncodeToString(EpochIDHashBytes(epochID))
}

func EpochIDHashScalar(epochID TypeEpochID) kyber.Scalar {
	return crypto.KyberSuite.Scalar().SetBytes(EpochIDHashBytes(epochID))
}

func EpochIDHashPoint(epochID TypeEpochID) kyber.Point {
	return crypto.KyberSuite.Point().Mul(EpochIDHashScalar(epochID), nil)
}
