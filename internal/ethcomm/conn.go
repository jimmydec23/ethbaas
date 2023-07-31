package ethcomm

import (
	"ethbaas/internal/config"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ChainConn struct {
	RpcCli  *rpc.Client
	GethCli *gethclient.Client
	EthCli  *ethclient.Client
}

func (c *ChainConn) Close() {
	c.RpcCli.Close()
}
func NewConn(port int32) (*ChainConn, error) {
	url := fmt.Sprintf("%s:%d", config.C.GetString("ethurl"), port)
	rpcCli, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	ethCli := ethclient.NewClient(rpcCli)
	gethCli := gethclient.New(rpcCli)
	c := &ChainConn{
		RpcCli:  rpcCli,
		EthCli:  ethCli,
		GethCli: gethCli,
	}
	return c, nil
}
