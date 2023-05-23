package chaincode

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/auti-project/auti-core/digest"
// 	"github.com/hyperledger/fabric-chaincode-go/shim"
// 	"github.com/hyperledger/fabric-contract-api-go/contractapi"
// 	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
// )

// // SmartContract provides functions for managing an Transaction.
// type SmartContract struct {
// 	contractapi.Contract
// }

// // InitLedger adds a base set of digests to the ledger.
// func (s *SmartContract) InitLedger(tci contractapi.TransactionContextInterface) error {
// 	digests := []digest.Digest{
// 		{
// 			Data:  "000",
// 			OrgID: "000",
// 		},
// 	}

// 	for _, d := range digests {
// 		key, val, err := d.KeyVal()
// 		if err != nil {
// 			return err
// 		}
// 		if err := tci.GetStub().PutState(key, val); err != nil {
// 			return fmt.Errorf("failed to put to world state: %v", err)
// 		}
// 	}
// 	return nil
// }

// // CreateTX issues a new transaction to the world state with given details.
// func (s *SmartContract) CreateTX(ctx contractapi.TransactionContextInterface,
// 	data, orgID string) (string, error) {
// 	d := digest.Digest{
// 		Data:  data,
// 		OrgID: orgID,
// 	}
// 	key, val, err := d.KeyVal()
// 	if err != nil {
// 		return "", err
// 	}
// 	var exists bool
// 	exists, err = s.TXExists(ctx, key)
// 	if err != nil {
// 		return "", err
// 	}
// 	if exists {
// 		return "", fmt.Errorf("the transaction %s already exists", key)
// 	}

// 	return key, ctx.GetStub().PutState(key, val)
// }

// func (s *SmartContract) CreateBatchTXs(ctx contractapi.TransactionContextInterface, digestListJSONString string) ([]string, error) {
// 	digestListJSONBytes, err := hex.DecodeString(digestListJSONString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var digestList []*digest.Digest
// 	err = json.Unmarshal(digestListJSONBytes, &digestList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	keys := make([]string, len(digestList))
// 	for i, tx := range digestList {
// 		key, val, err := tx.KeyVal()
// 		if err != nil {
// 			return nil, err
// 		}
// 		var exists bool
// 		exists, err = s.TXExists(ctx, key)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if exists {
// 			return nil, fmt.Errorf("the transaction %s already exists", key)
// 		}
// 		if err := ctx.GetStub().PutState(key, val); err != nil {
// 			return nil, fmt.Errorf("failed to put to world state: %v", err)
// 		}
// 		keys[i] = key
// 	}
// 	return keys, nil
// }

// // ReadTX returns the transaction stored in the world state with given id.
// func (s *SmartContract) ReadTX(ctx contractapi.TransactionContextInterface, key string) (*digest.Digest, error) {
// 	tx, err := ctx.GetStub().GetState(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read from world state: %v", err)
// 	}
// 	if tx == nil {
// 		return nil, fmt.Errorf("the transaction %s does not exist", key)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	var digestObj digest.Digest
// 	err = json.Unmarshal(tx, &digestObj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &digestObj, nil
// }

// // DeleteTX deletes a given transaction from the world state.
// func (s *SmartContract) DeleteTX(ctx contractapi.TransactionContextInterface, key string) error {
// 	exists, err := s.TXExists(ctx, key)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the transaction %s does not exist", key)
// 	}
// 	return ctx.GetStub().DelState(key)
// }

// // TXExists returns true when transaction with given ID exists in world state.
// func (s *SmartContract) TXExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
// 	transactionJSON, err := ctx.GetStub().GetState(key)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to read from world state: %v", err)
// 	}
// 	return transactionJSON != nil, nil
// }

// // GetAllTXs returns all transactions found in world state.
// func (s *SmartContract) GetAllTXs(ctx contractapi.TransactionContextInterface) (txList []*digest.Digest, err error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all transactions in the chaincode namespace.
// 	var iter shim.StateQueryIteratorInterface
// 	iter, err = ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return
// 	}
// 	defer func(resultsIterator shim.StateQueryIteratorInterface) {
// 		err = resultsIterator.Close()
// 	}(iter)

// 	var response *queryresult.KV
// 	for iter.HasNext() {
// 		response, err = iter.Next()
// 		if err != nil {
// 			return
// 		}
// 		var d digest.Digest
// 		err = json.Unmarshal(response.Value, &d)
// 		if err != nil {
// 			return
// 		}
// 		txList = append(txList, &d)
// 	}
// 	return
// }