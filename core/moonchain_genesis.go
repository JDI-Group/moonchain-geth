package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	moonchainGenesis "github.com/ethereum/go-ethereum/core/moonchain_genesis"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

var (
	MoonchainMainnetOntakeBlock = new(big.Int).SetUint64(0)
	MoonchainHudsonOntakeBlock  = new(big.Int).SetUint64(0)
)

// MoonchainGenesisBlock returns the Moonchain network genesis block configs.
func MoonchainGenesisBlock(networkID uint64) *Genesis {
	chainConfig := params.MoonchainChainConfig

	var allocJSON []byte
	switch networkID {
	case params.MoonchainMainnetNetworkID.Uint64():
		chainConfig.ChainID = params.MoonchainMainnetNetworkID
		chainConfig.OntakeBlock = MoonchainMainnetOntakeBlock
		allocJSON = moonchainGenesis.MoonchainMainnetGenesisAllocJSON
	case params.MoonchainHudsonNetworkID.Uint64():
		chainConfig.ChainID = params.MoonchainHudsonNetworkID
		chainConfig.OntakeBlock = MoonchainHudsonOntakeBlock
		allocJSON = moonchainGenesis.MoonchainHudsonGenesisAllocJSON
	default:
		chainConfig.ChainID = params.MoonchainHudsonNetworkID
		chainConfig.OntakeBlock = MoonchainHudsonOntakeBlock
		allocJSON = moonchainGenesis.MoonchainHudsonGenesisAllocJSON
	}

	var alloc GenesisAlloc
	if err := alloc.UnmarshalJSON(allocJSON); err != nil {
		log.Crit("unmarshal alloc json error", "error", err)
	}

	return &Genesis{
		Config:     chainConfig,
		ExtraData:  []byte{},
		GasLimit:   uint64(6000000),
		Difficulty: common.Big0,
		Alloc:      alloc,
		GasUsed:    0,
		BaseFee:    new(big.Int).SetUint64(10000000),
	}
}
