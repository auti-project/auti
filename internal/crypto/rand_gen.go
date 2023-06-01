package crypto

import (
	"crypto/rand"

	"github.com/auti-project/auti/internal/constants"
	"go.dedis.ch/kyber/v3"
)

func RandBytes() ([]byte, error) {
	result := make([]byte, constants.SecurityParameterBytes)
	_, err := rand.Read(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func RandScalars(size uint) []kyber.Scalar {
	results := make([]kyber.Scalar, size)
	for i := uint(0); i < size; i++ {
		randScalar := KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
		results[i] = randScalar
	}
	return results
}
