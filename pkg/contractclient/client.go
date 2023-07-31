package contractclient

import (
	"context"
	"ethbaas/internal/db"
	"ethbaas/internal/ethcomm"
	"ethbaas/internal/model"
	"ethbaas/pkg/projclient"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client struct {
	db      *db.Client
	projCli *projclient.Client
}

func NewClient(db *db.Client) *Client {
	c := &Client{
		db:      db,
		projCli: projclient.NewClient(db),
	}
	return c
}

func (c *Client) Deploy(projName string, contract *model.Contract) (*common.Address, error) {
	dbproj, err := c.db.GetProject(projName)
	if err != nil {
		return nil, err
	}

	apiPort := dbproj.Str2Port()[0]

	conn, err := ethcomm.NewConn(apiPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	opts, err := ethcomm.GenTxOpts(conn)
	if err != nil {
		return nil, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return nil, err
	}

	contractBin := common.FromHex(contract.BIN)

	version := "1.0"
	address, tx, _, err := bind.DeployContract(opts, contractAbi, contractBin, conn.EthCli, version)
	if err != nil {
		return nil, err
	}
	fmt.Println("Contract deployed at:", address.Hex())
	fmt.Println("Contract deploy tx:", tx.Hash().Hex())
	err = c.db.AddContract(&db.Contract{
		Name:    contract.Name,
		Proj:    dbproj.Name,
		Created: time.Now().Unix(),
		Address: address.Hex(),
		ABI:     contract.ABI,
		BIN:     contract.BIN,
	})

	return &address, nil
}

func (c *Client) List() (int64, []db.Contract, error) {
	return c.db.ListContract()
}

func (c *Client) Query(contractName, methodName, inputs string) (interface{}, error) {
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

	contractAbi, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return nil, err
	}

	inputList := []interface{}{}
	method, ok := contractAbi.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("Method not found.")
	}

	userInputs := strings.Split(inputs, ",")
	for i := 0; i < len(method.Inputs); i++ {
		methodInput := method.Inputs[i]
		userInput := userInputs[i]
		switch methodInput.Type.GetType() {
		case reflect.TypeOf([32]uint8{}):
			transInput := [32]uint8{}
			copy(transInput[:], []byte(userInput))
			inputList = append(inputList, transInput)
		default:
			inputList = append(inputList, userInput)
		}
	}

	data, err := contractAbi.Pack(methodName, inputList...)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}
	return conn.EthCli.CallContract(context.Background(), msg, nil)
}

func (c *Client) Write(contractName, methodName, inputs string) (string, error) {
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
		return "", nil
	}
	defer conn.Close()

	address := common.HexToAddress(contract.Address)

	contractAbi, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return "", err
	}

	inputList := []interface{}{}
	method, ok := contractAbi.Methods[methodName]
	if !ok {
		return "", fmt.Errorf("Method not found.")
	}

	userInputs := strings.Split(inputs, ",")
	for i := 0; i < len(method.Inputs); i++ {
		methodInput := method.Inputs[i]
		userInput := userInputs[i]
		switch methodInput.Type.GetType() {
		case reflect.TypeOf([32]uint8{}):
			transInput := [32]uint8{}
			copy(transInput[:], []byte(userInput))
			inputList = append(inputList, transInput)
		default:
			inputList = append(inputList, userInput)
		}
	}

	data, err := contractAbi.Pack(methodName, inputList...)
	if err != nil {
		return "", err
	}

	txOpt, err := ethcomm.GenTxOpts(conn)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(
		txOpt.Nonce.Uint64(),
		address,
		big.NewInt(0),
		txOpt.GasLimit,
		txOpt.GasPrice,
		data,
	)
	chainId, err := conn.EthCli.ChainID(context.Background())
	if err != nil {
		return "", err
	}

	privateKey, err := ethcomm.GetAdminPk()
	if err != nil {
		return "", err
	}

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return "", err
	}

	err = conn.EthCli.SendTransaction(context.Background(), signTx)
	if err != nil {
		return "", err
	}

	return signTx.Hash().Hex(), nil
}
