package core

func Commit(amount int64) ([]byte, error) {
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, err
	}
	amountScalar := kyberSuite.Scalar().SetBytes(amountBytes)
	commitment := kyberSuite.Point().Mul(amountScalar, G)
	randScalar := kyberSuite.Scalar().Pick(kyberSuite.RandomStream())
	randPoint := kyberSuite.Point().Mul(randScalar, H)
	commitment.Add(commitment, randPoint)
	return commitment.MarshalBinary()
}
