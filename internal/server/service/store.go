package service

import (
	"bytes"
	"crypto/ecdsa"
	"ethbaas/contract/store"
	"ethbaas/internal/db"
	"ethbaas/internal/ethcomm"
	"ethbaas/internal/log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type StoreSvc struct {
	dbClient *db.Client
	conn     *ethcomm.ChainConn
	instance *store.Store
	adminPk  *ecdsa.PrivateKey
	mu       sync.Mutex
}

func NewStoreSvc(dbClient *db.Client) *StoreSvc {
	s := &StoreSvc{
		dbClient: dbClient,
		mu:       sync.Mutex{},
	}
	return s
}

func (s *StoreSvc) setup() error {
	contract, err := s.dbClient.GetContract("store")
	if err != nil {
		return err
	}

	dbproj, err := s.dbClient.GetProject(contract.Proj)
	if err != nil {
		return err
	}

	apiPort := dbproj.Str2Port()[0]

	conn, err := ethcomm.NewConn(apiPort)
	if err != nil {
		return err
	}

	address := common.HexToAddress(contract.Address)
	instance, err := store.NewStore(address, conn.EthCli)
	if err != nil {
		return err
	}

	adminPk, err := ethcomm.GetAdminPk()
	if err != nil {
		return err
	}

	s.adminPk = adminPk
	s.conn = conn
	s.instance = instance
	return nil
}

func (s *StoreSvc) Query(key string) (string, error) {
	if s.instance == nil {
		err := s.setup()
		if err != nil {
			return "", err
		}
	}
	bkey := [32]byte{}
	copy(bkey[:], []byte(key))
	value, err := s.instance.Items(nil, bkey)
	if err != nil {
		return "", err
	}
	v := string(bytes.Trim(value[:], "\u0000"))
	return v, nil
}

func (s *StoreSvc) Write(key, value string) (string, error) {
	if s.instance == nil {
		err := s.setup()
		if err != nil {
			return "", err
		}
	}
	// tx nonce must not equl
	s.mu.Lock()
	defer s.mu.Unlock()

	auth, err := ethcomm.GenTxOpts(s.conn)
	log.Logger.Info("tx nonce:", auth.Nonce)
	if err != nil {
		return "", err
	}

	bkey := [32]byte{}
	bvalue := [32]byte{}
	copy(bkey[:], []byte(key))
	copy(bvalue[:], []byte(value))
	tx, err := s.instance.SetItem(auth, bkey, bvalue)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}
