package utils

import (
	"os"

	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli/v2"
)

var (
	MoonchainFlag = cli.BoolFlag{
		Name:  "moonchain",
		Usage: "Moonchain network",
	}
)

// RegisterMoonchainAPIs initializes and registers the Moonchain RPC APIs.
func RegisterMoonchainAPIs(stack *node.Node, cfg *ethconfig.Config, backend *eth.Ethereum) {
	if os.Getenv("MOONCHAIN_TEST") != "" {
		return
	}
	// Add methods under "moonchain_" RPC namespace to the available APIs list
	stack.RegisterAPIs([]rpc.API{
		{
			Namespace: "moonchain",
			Version:   params.VersionWithMeta,
			Service:   eth.NewMoonchainAPIBackend(backend),
			Public:    true,
		},
		{
			Namespace:     "moonchainAuth",
			Version:       params.VersionWithMeta,
			Service:       eth.NewMoonchainAuthAPIBackend(backend),
			Authenticated: true,
		},
	})
}
