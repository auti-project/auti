package orgchain

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/auti-project/auti/internal/transaction"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	channelName    = "mychannel"
	contractType   = "auti-org-chain"
	orgWalletPath  = "orgWallet"
	orgWalletLabel = "orgAPPUser"
	audWalletPath  = "audWallet"
	audWalletLabel = "audAPPUser"
	org1MSPid      = "Org1MSP"
	aud1MSPid      = "Aud1MSP"
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

	org1CCPPath  string
	org1CREDPath string
	org1CertPath string
	org1KeyDir   string

	aud1CCPPath  string
	aud1CREDPath string
	aud1CertPath string
	aud1KeyDir   string
)

func init() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}
	err = os.RemoveAll(orgWalletPath)
	if err != nil {
		log.Fatalf("Error removing wallet directory: %v", err)
	}
	err = os.RemoveAll(audWalletPath)
	if err != nil {
		log.Fatalf("Error removing wallet directory: %v", err)
	}

	autiLocalChainDir := os.Getenv("AUTI_LOCAL_CHAIN_DIR")

	fabloFilePath = filepath.Join(autiLocalChainDir, "fablo-target", "fabric-config")

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

	aud1CCPPath = filepath.Join(fabloFilePath, "connection-profiles", "connection-profile-aud1.yaml")
	aud1CREDPath = filepath.Join(
		fabloFilePath,
		"crypto-config",
		"peerOrganizations",
		"aud1.example.com",
		"users",
		"User1@aud1.example.com",
		"msp",
	)
	aud1CertPath = filepath.Join(aud1CREDPath, "signcerts", "User1@aud1.example.com-cert.pem")
	aud1KeyDir = filepath.Join(aud1CREDPath, "keystore")
}

type Controller struct {
	gw *gateway.Gateway
	ct *gateway.Contract
}

// NewController starts a new service instance
func NewController(walletPath, walletLabel, ccpPath string) (*Controller, error) {
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		return nil, err
	}
	if !wallet.Exists(walletLabel) {
		// TODO: bad practice, should be removed
		if walletLabel == orgWalletLabel {
			if err = populateOrgWallet(wallet); err != nil {
				return nil, err
			}
		} else {
			if err = populateAudWallet(wallet); err != nil {
				return nil, err
			}
		}
	}
	var gw *gateway.Gateway
	if gw, err = gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, walletLabel),
	); err != nil {
		return nil, err
	}
	network, err := gw.GetNetwork(channelName)
	if err != nil {
		return nil, err
	}
	contract := network.GetContract(contractType)
	return &Controller{gw: gw, ct: contract}, nil
}

func (s *Controller) Close() {
	s.gw.Close()
}

func (s *Controller) SubmitTX(tx *transaction.CLOLCLocalOnChain) (string, error) {
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

func (s *Controller) GetAllTXs() ([]*transaction.CLOLCLocalOnChain, error) {
	results, err := s.ct.EvaluateTransaction(getAllTXFuncName)
	if err != nil {
		return nil, err
	}
	var txList []*transaction.CLOLCLocalOnChain
	err = json.Unmarshal(results, &txList)
	if err != nil {
		return nil, err
	}
	return txList, nil
}

func (s *Controller) ReadTX(id string) (*transaction.CLOLCLocalOnChain, error) {
	result, err := s.ct.EvaluateTransaction(readTXFuncName, id)
	if err != nil {
		return nil, err
	}
	var tx transaction.CLOLCLocalOnChain
	err = json.Unmarshal(result, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *Controller) SubmitBatchTXs(txList []*transaction.CLOLCLocalOnChain) ([]string, error) {
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

func populateOrgWallet(wallet *gateway.Wallet) error {
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

	return wallet.Put(orgWalletLabel, identity)
}

func populateAudWallet(wallet *gateway.Wallet) error {
	cert, err := os.ReadFile(filepath.Clean(aud1CertPath))
	if err != nil {
		return err
	}

	files, err := os.ReadDir(aud1KeyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(aud1KeyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(aud1MSPid, string(cert), string(key))

	return wallet.Put(audWalletLabel, identity)
}

const (
	txThreshold = 10000
	txIDLogPath = "lc_tx_id.log"
)

func BenchReadTX() error {
	f, err := os.Open(txIDLogPath)
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	var txIDList []string
	for fileScanner.Scan() {
		txIDList = append(txIDList, fileScanner.Text())
	}
	err = f.Close()
	if err != nil {
		return err
	}
	lc, err := NewController(audWalletPath, audWalletLabel, aud1CCPPath)
	if err != nil {
		return err
	}
	defer lc.Close()
	idx := rand.Int() % len(txIDList)
	_, err = lc.ReadTX(txIDList[idx])
	return err
}

func ReadAllTXs() error {
	lc, err := NewController(audWalletPath, audWalletLabel, aud1CCPPath)
	if err != nil {
		return err
	}
	defer lc.Close()
	_, err = lc.GetAllTXs()
	return err
}

func SubmitTX(numTXs int) ([]string, error) {
	lc, err := NewController(orgWalletPath, orgWalletLabel, org1CCPPath)
	if err != nil {
		return nil, err
	}
	defer lc.Close()
	dummyTXs := DummyOnChainTransactions(numTXs)
	var txIDs []string
	for j := 0; j < numTXs; j += txThreshold {
		right := j + txThreshold
		if right > numTXs {
			right = numTXs
		}
		batchTXIDs, err := lc.SubmitBatchTXs(dummyTXs[j:right])
		if err != nil {
			return nil, err
		}
		txIDs = append(txIDs, batchTXIDs...)
	}
	return txIDs, nil
}

func SaveTXIDs(txIDs []string) error {
	f, err := os.OpenFile(txIDLogPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	for _, id := range txIDs {
		if _, err = f.WriteString(id + "\n"); err != nil {
			return err
		}
	}
	return nil
}
