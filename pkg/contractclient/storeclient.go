package contractclient

import (
	"ethbaas/contract/store"
	"ethbaas/internal/db"
	"ethbaas/internal/ethcomm"
	"ethbaas/pkg/projclient"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type StoreClient struct {
	db      *db.Client
	projCli *projclient.Client
}

func NewStoreClient(db *db.Client) *StoreClient {
	c := &StoreClient{
		db:      db,
		projCli: projclient.NewClient(db),
	}
	return c
}

// deploy kv store contract
func (c *StoreClient) Deploy(projName, contractName, contractPath string) error {
	dbproj, err := c.projCli.Get(projName)
	if err != nil {
		return err
	}

	apiPort := dbproj.Str2Port()[0]

	conn, err := ethcomm.NewConn(apiPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	auth, err := ethcomm.GenTxOpts(conn)
	if err != nil {
		return err
	}

	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, conn.EthCli, input)
	if err != nil {
		return err
	}

	fmt.Println("Contract deployed at:", address.Hex())
	fmt.Println("Contract deploy tx:", tx.Hash().Hex())
	_ = instance

	err = c.db.AddContract(&db.Contract{
		Name:    contractName,
		Proj:    dbproj.Name,
		Created: time.Now().Unix(),
		Address: address.Hex(),
		ABI:     store.StoreABI,
		BIN:     store.StoreBin,
	})
	return err
}

// query kv store contract
func (c *StoreClient) Query(contractName, key string) (interface{}, error) {
	contract, err := c.db.GetContract(contractName)
	if err != nil {
		return nil, err
	}

	dbproj, err := c.db.GetProject(contract.Proj)
	if err != nil {
		return nil, err
	}

	apiPort := dbproj.Str2Port()[0]

	conn, err := ethcomm.NewConn(apiPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	address := common.HexToAddress(contract.Address)
	instance, err := store.NewStore(address, conn.EthCli)
	if err != nil {
		return nil, err
	}

	bkey := [32]byte{}
	copy(bkey[:], []byte(key))
	value, err := instance.Items(nil, bkey)
	if err != nil {
		return nil, err
	}
	return string(value[:]), nil
}

// write to kv store contract
func (c *StoreClient) Write(contractName, key, value string) (string, error) {
	contract, err := c.db.GetContract(contractName)
	if err != nil {
		return "", err
	}

	dbproj, err := c.db.GetProject(contract.Proj)
	if err != nil {
		return "", err
	}

	apiPort := dbproj.Str2Port()[0]

	conn, err := ethcomm.NewConn(apiPort)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	address := common.HexToAddress(contract.Address)
	instance, err := store.NewStore(address, conn.EthCli)
	if err != nil {
		return "", err
	}
	auth, err := ethcomm.GenTxOpts(conn)
	if err != nil {
		return "", err
	}

	bkey := [32]byte{}
	bvalue := [32]byte{}
	copy(bkey[:], []byte(key))
	copy(bvalue[:], []byte(value))
	tx, err := instance.SetItem(auth, bkey, bvalue)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}
