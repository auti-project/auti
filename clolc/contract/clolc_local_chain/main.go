package main

import (
	"log"

	"github.com/auti-project/auti/clolc/contract/clolc_local_chain/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	sc := new(chaincode.SmartContract)
	cc, err := contractapi.NewChaincode(sc)
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
