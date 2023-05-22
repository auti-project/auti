package crypto

import (
	"crypto/rand"
	"math/big"

	"github.com/auti-project/auti/core/constants"
)

var (
	big1  = big.NewInt(1)
	limit = new(big.Int).Lsh(big1, constants.SecurityParameter)
)

func RandInt() (*big.Int, error) {
	return rand.Int(rand.Reader, limit)
}

func RandIntList(size uint) ([]*big.Int, error) {
	results := make([]*big.Int, size)
	for i := uint(0); i < size; i++ {
		randInt, err := RandInt()
		if err != nil {
			return nil, err
		}
		results[i] = randInt
	}
	return results, nil
}
