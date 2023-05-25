package ledger

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/auti-project/auti/internal/transaction"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	channelName  = "mychannel"
	contractType = "auti-local-chain"
	walletPath   = "wallet"
	walletLabel  = "appUser"
	org1MSPid    = "Org1MSP"
)

const (
	createTXFuncName      = "CreateTX"
	createBatchTXFuncName = "CreateBatchTXs"
	txExistsName          = "TXExists"
	getAllTXFuncName      = "GetAllTXs"
	readTXFuncName        = "ReadTX"
)

var (
	fabloFilePath string
	org1CCPPath   string
	org1CREDPath  string
	org1CertPath  string
	org1KeyDir    string
)

func init() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}
	err = os.RemoveAll(walletPath)
	if err != nil {
		log.Fatalf("Error removing wallet directory: %v", err)
	}
	autiOrgGlobalDir := os.Getenv("AUTI_ORG_GLOBAL_DIR")

	fabloFilePath = filepath.Join(autiOrgGlobalDir, "fablo-target", "fabric-config")
	org1CCPPath = filepath.Join(fabloFilePath, "connection-profiles", "connection-profile-org1.yaml")
	org1CREDPath = filepath.Join(
		fabloFilePath,
		"crypto-config",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)
	org1CertPath = filepath.Join(org1CREDPath, "signcerts", "User1@org1.example.com-cert.pem")
	org1KeyDir = filepath.Join(org1CREDPath, "keystore")
}

type Controller struct {
	gw *gateway.Gateway
	ct *gateway.Contract
}

// NewController starts a new service instance
func NewController() *Controller {
	service := new(Controller)
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(walletLabel) {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}
	var gw *gateway.Gateway
	gw, err = gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(org1CCPPath))),
		gateway.WithIdentity(wallet, walletLabel),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	service.gw = gw
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract := network.GetContract(contractType)
	service.ct = contract
	// log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	// result, err := contract.SubmitTransaction("InitLedger")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))
	return service
}

func (s *Controller) Close() {
	s.gw.Close()
}

func (s *Controller) SubmitTX(tx *transaction.CLOLCOnChain) (string, error) {
	// log.Println("--> Submit Transaction: Invoke, function that adds a new asset")
	txID, err := s.ct.SubmitTransaction(createTXFuncName,
		tx.CounterParty,
		tx.Commitment,
		tx.Timestamp,
	)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	return string(txID), nil
}

func (s *Controller) TXExists(txID string) (bool, error) {
	resBytes, err := s.ct.EvaluateTransaction(txExistsName, txID)
	if err != nil {
		return false, err
	}
	var result bool
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *Controller) GetAllTXs() ([]*transaction.CLOLCOnChain, error) {
	results, err := s.ct.EvaluateTransaction(getAllTXFuncName)
	if err != nil {
		return nil, err
	}
	var txList []*transaction.CLOLCOnChain
	err = json.Unmarshal(results, &txList)
	if err != nil {
		return nil, err
	}
	return txList, nil
}

func (s *Controller) ReadTX(id string) (*transaction.CLOLCOnChain, error) {
	result, err := s.ct.EvaluateTransaction(readTXFuncName, id)
	if err != nil {
		return nil, err
	}
	var tx transaction.CLOLCOnChain
	err = json.Unmarshal(result, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *Controller) SubmitBatchTXs(txList []*transaction.CLOLCOnChain) ([]string, error) {
	txListJSON, err := json.Marshal(txList)
	if err != nil {
		return nil, err
	}
	txListJSONstr := hex.EncodeToString(txListJSON)
	resBytes, err := s.ct.SubmitTransaction(createBatchTXFuncName, txListJSONstr)
	if err != nil {
		return nil, err
	}
	var txIDList []string
	err = json.Unmarshal(resBytes, &txIDList)
	if err != nil {
		return nil, err
	}
	return txIDList, nil
}

func populateWallet(wallet *gateway.Wallet) error {
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(org1CertPath))
	if err != nil {
		return err
	}

	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(org1KeyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(org1KeyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(org1MSPid, string(cert), string(key))

	return wallet.Put(walletLabel, identity)
}
