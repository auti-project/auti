package organization

import "math/big"

type TypeID string

type Organization struct {
	ID          TypeID
	epochRandID *big.Int
}

func New(id string) *Organization {
	return &Organization{
		ID: TypeID(id),
	}
}

func (o *Organization) SetEpochID(randID *big.Int) {
	o.epochRandID = randID
}
