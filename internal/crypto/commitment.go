package crypto

import "go.dedis.ch/kyber/v3"

func PedersenCommit(amount int64) (kyber.Point, error) {
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, err
	}
	amountScalar := kyberSuite.Scalar().SetBytes(amountBytes)
	commitment := kyberSuite.Point().Mul(amountScalar, pointG)
	randScalar := kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	randPoint := kyberSuite.Point().Mul(randScalar, pointH)
	commitment.Add(commitment, randPoint)
	return commitment, nil
}
