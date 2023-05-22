package crypto

func PedersenCommit(amount int64) ([]byte, error) {
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, err
	}
	amountScalar := kyberSuite.Scalar().SetBytes(amountBytes)
	commitment := kyberSuite.Point().Mul(amountScalar, pointG)
	randScalar := kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	randPoint := kyberSuite.Point().Mul(randScalar, pointH)
	commitment.Add(commitment, randPoint)
	return commitment.MarshalBinary()
}
