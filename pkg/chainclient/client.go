package chainclient

import (
	"context"
	"ethbaas/internal/db"
	"ethbaas/internal/ethcomm"
	"ethbaas/internal/k8s"
	"ethbaas/internal/model"
	"ethbaas/pkg/projclient"
	"fmt"
	"math/big"
	"strings"
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

func (c *Client) Info(projName string) ([]ChainInfo, error) {
	dbproj, err := c.projCli.Get(projName)
	if err != nil {
		return nil, err
	}

	ports := dbproj.Str2Port()
	infoes := []ChainInfo{}

	for _, port := range ports {
		conn, err := ethcomm.NewConn(port)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		networkId, err := conn.EthCli.NetworkID(context.Background())
		if err != nil {
			return nil, err
		}
		blockNum, err := conn.EthCli.BlockNumber(context.Background())
		if err != nil {
			return nil, err
		}
		pc, err := conn.EthCli.PeerCount(context.Background())
		if err != nil {
			return nil, err
		}

		sp, err := conn.EthCli.SyncProgress(context.Background())
		if err != nil {
			return nil, err
		}

		current := blockNum
		highest := blockNum
		if sp != nil {
			current = sp.CurrentBlock
			highest = sp.HighestBlock
		}
		block, err := conn.EthCli.BlockByNumber(context.Background(), big.NewInt(int64(current)))
		if err != nil {
			return nil, err
		}

		conn.GethCli.GetNodeInfo(context.Background())
		node, err := conn.GethCli.GetNodeInfo(context.Background())
		if err != nil {
			return nil, err
		}

		enodeFirst := strings.Split(node.Enode, "@")[0][:20] + "..."
		info := ChainInfo{
			NetworkID: networkId,
			Current:   current,
			Highest:   highest,
			Syncing:   sp != nil,
			PeerCount: pc,
			ENode:     enodeFirst,
			Diff:      block.Difficulty(),
		}
		infoes = append(infoes, info)
	}
	return infoes, nil
}

func (c *Client) Pods(p *model.Project) ([]k8s.Pod, error) {
	return k8s.GetPods(p.NS())
}

func (c *Client) Cluster(p *model.Project) error {
	dbproj, err := c.db.GetProject(p.Name)
	if err != nil {
		return err
	}
	ports := dbproj.Str2Port()

	node0Conn, err := ethcomm.NewConn(ports[0])
	if err != nil {
		return nil
	}
	defer node0Conn.Close()

	for i := 1; i < len(ports); i++ {
		apiPort := ports[i]
		conn, err := ethcomm.NewConn(apiPort)
		if err != nil {
			return err
		}
		defer conn.Close()

		node, err := conn.GethCli.GetNodeInfo(context.Background())
		if err != nil {
			return err
		}

		enode := node.Enode
		enode = strings.ReplaceAll(enode, "127.0.0.1", fmt.Sprintf("node%d", i))
		addResult := false

		if err := node0Conn.RpcCli.Call(&addResult, "admin_addPeer", enode); err != nil {
			return err
		}
	}
	return nil
}
