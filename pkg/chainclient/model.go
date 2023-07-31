package chainclient

import "math/big"

type ChainInfo struct {
	NetworkID *big.Int
	Current   uint64
	Highest   uint64
	Syncing   bool
	PeerCount uint64
	ENode     string
	Diff      *big.Int
}
