package closc

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/auti-project/auti/internal/organization"
)

var sha256Func = sha256.New()

type Organization struct {
	ID     organization.TypeID
	IDHash string
}

func New(id string) *Organization {
	defer sha256Func.Reset()
	sha256Func.Write([]byte(id))
	idHash := hex.EncodeToString(sha256Func.Sum(nil))
	org := &Organization{
		ID:     organization.TypeID(id),
		IDHash: idHash,
	}
	return org
}
