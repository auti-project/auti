package crypto

import "go.dedis.ch/kyber/v3"

func PedersenCommit(amount int64) (kyber.Point, kyber.Scalar, error) {
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, nil, err
	}
	amountScalar := KyberSuite.Scalar().SetBytes(amountBytes)
	commitment := KyberSuite.Point().Mul(amountScalar, PointG)
	randScalar := KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
	randPoint := KyberSuite.Point().Mul(randScalar, PointH)
	commitment.Add(commitment, randPoint)
	return commitment, randScalar, nil
}
