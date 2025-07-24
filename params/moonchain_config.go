package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Network IDs
var (
	MoonchainMainnetNetworkID = big.NewInt(999888)
	MoonchainHudsonNetworkID  = big.NewInt(177888)
)

var networkIDToChainConfigByMoonchain = map[*big.Int]*ChainConfig{
	MoonchainMainnetNetworkID:  MoonchainChainConfig,
	MoonchainHudsonNetworkID:   MoonchainChainConfig,
	MainnetChainConfig.ChainID: MainnetChainConfig,
	SepoliaChainConfig.ChainID: SepoliaChainConfig,
	TestChainConfig.ChainID:    TestChainConfig,
	NonActivatedConfig.ChainID: NonActivatedConfig,
}

func NetworkIDToChainConfigOrDefaultByMoonchain(networkID *big.Int) *ChainConfig {
	if config, ok := networkIDToChainConfigByMoonchain[networkID]; ok {
		return config
	}

	return AllEthashProtocolChanges
}

var MoonchainChainConfig = &ChainConfig{
	ChainID:                       MoonchainHudsonNetworkID, // Use Moonchain Hudson network ID by default.
	HomesteadBlock:                common.Big0,
	EIP150Block:                   common.Big0,
	EIP155Block:                   common.Big0,
	EIP158Block:                   common.Big0,
	ByzantiumBlock:                common.Big0,
	ConstantinopleBlock:           common.Big0,
	PetersburgBlock:               common.Big0,
	IstanbulBlock:                 common.Big0,
	BerlinBlock:                   common.Big0,
	LondonBlock:                   common.Big0,
	ShanghaiTime:                  u64(0),
	MergeNetsplitBlock:            nil,
	TerminalTotalDifficulty:       common.Big0,
	TerminalTotalDifficultyPassed: true,
	Moonchain:                     true,
}
