package core

import (
	"crypto/rand"
	"math/big"
)

func RandInt() (*big.Int, error) {
	return rand.Int(rand.Reader, RandIntLimit)
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
