package crypto

import "go.dedis.ch/kyber/v3"

func PedersenCommit(amount int64) (kyber.Point, kyber.Scalar, error) {
	amountScalar, err := amountToScalar(amount)
	if err != nil {
		return nil, nil, err
	}
	commitment := KyberSuite.Point().Mul(amountScalar, PointG)
	randScalar := KyberSuite.Scalar().Pick(KyberSuite.RandomStream())
	randPoint := KyberSuite.Point().Mul(randScalar, PointH)
	commitment.Add(commitment, randPoint)
	return commitment, randScalar, nil
}

func amountToScalar(amount int64) (kyber.Scalar, error) {
	positive := true
	if amount < 0 {
		amount = -amount
		positive = false
	}
	amountBytes, err := int64ToBytes(amount)
	if err != nil {
		return nil, err
	}
	amountScalar := KyberSuite.Scalar().SetBytes(amountBytes)
	if !positive {
		amountScalar.Neg(amountScalar)
	}
	return amountScalar, nil
}
