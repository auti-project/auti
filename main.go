package main

import (
	"fmt"

	"github.com/auti-project/auti/internal/core"
)

func main() {
	//pubKey, privKey, _ := core.KeyGen()
	//a := big.NewInt(12)
	//cipherText, err := core.Encrypt(pubKey, a)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(cipherText)
	//plainText := core.Decrypt(privKey, cipherText)
	//fmt.Println(plainText)
	fmt.Println(core.H.MarshalBinary())
}
