package organization

import (
	"crypto/sha256"
	"encoding/hex"

	"go.dedis.ch/kyber/v3"

	"github.com/auti-project/auti/internal/crypto"
)

var (
	Sha256Func = sha256.New()
)

type TypeID string
type TypeEpochID []byte

func IDHashBytes(id TypeID) []byte {
	defer Sha256Func.Reset()
	Sha256Func.Write([]byte(id))
	return Sha256Func.Sum(nil)
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

func IDHashKey(orgIDHash1, orgIDHash2 string) [2]string {
	if orgIDHash1 < orgIDHash2 {
		return [2]string{orgIDHash1, orgIDHash2}
	}
	return [2]string{orgIDHash2, orgIDHash1}
}

func EpochIDHashBytes(epochID TypeEpochID) []byte {
	defer Sha256Func.Reset()
	Sha256Func.Write(epochID)
	return Sha256Func.Sum(nil)
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
