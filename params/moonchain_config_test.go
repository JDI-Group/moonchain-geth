package params

import (
	"math/big"
	"testing"
)

func TestNetworkIDToChainConfigOrDefaultByMoonchain(t *testing.T) {
	tests := []struct {
		name            string
		networkID       *big.Int
		wantChainConfig *ChainConfig
	}{
		{
			"moonchain",
			MoonchainMainnetNetworkID,
			MoonchainChainConfig,
		},
		{
			"moonchain-hudson",
			MoonchainHudsonNetworkID,
			MoonchainChainConfig,
		},
		{
			"mainnet",
			MainnetChainConfig.ChainID,
			MainnetChainConfig,
		},
		{
			"sepolia",
			SepoliaChainConfig.ChainID,
			SepoliaChainConfig,
		},
		{
			"doesntExist",
			big.NewInt(89390218390),
			AllEthashProtocolChanges,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if config := NetworkIDToChainConfigOrDefaultByMoonchain(tt.networkID); config != tt.wantChainConfig {
				t.Fatalf("expected %v, got %v", config, tt.wantChainConfig)
			}
		})
	}
}
