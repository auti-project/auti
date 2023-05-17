package core

import (
	"math/big"
	"testing"
)

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		amount  *big.Int
		wantErr bool
	}{
		{
			name:    "test_1",
			amount:  big.NewInt(1),
			wantErr: false,
		},
		{
			name:    "test_11",
			amount:  big.NewInt(11),
			wantErr: false,
		},
		{
			name:    "test_123123123123",
			amount:  big.NewInt(123123123123),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, publicKey, err := KeyGen()
			if err != nil {
				t.Errorf("KeyGen() error = %v", err)
				return
			}
			var cipherText *CipherText
			cipherText, err = Encrypt(publicKey, tt.amount)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}
			var plainText *big.Int
			plainText, err = Decrypt(privateKey, cipherText)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if plainText.Cmp(tt.amount) != 0 {
				t.Errorf("Decrypt() plainText = %v, amount %v", plainText, tt.amount)
			}
		})
	}
}
